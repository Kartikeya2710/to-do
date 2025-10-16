package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handlers) AddTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Bad reqeust", http.StatusBadRequest)

		return
	}

	dbCtx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	insertResult, err := h.collection.InsertOne(dbCtx, task)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error"), http.StatusInternalServerError)

		return
	}

	insertedId := insertResult.InsertedID.(bson.ObjectID)
	if err := json.NewEncoder(w).Encode(insertedId); err != nil {
		http.Error(w, fmt.Sprintf("Internal server error"), http.StatusInternalServerError)
	}
}

package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Kartikeya2710/to-do/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handlers) AddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Bad reqeust", http.StatusBadRequest)
		h.logger.Fatalf("Error decoding task from request body: %v\n", err)

		return
	}

	dbCtx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	insertResult, err := h.collection.InsertOne(dbCtx, task)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Fatalf("Error inserting document into collection: %v\n", err)

		return
	}

	insertedId := insertResult.InsertedID.(bson.ObjectID)
	if err := json.NewEncoder(w).Encode(insertedId); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Fatalf("Error sending inserted document id in response: %v\n", err)

		return
	}
}

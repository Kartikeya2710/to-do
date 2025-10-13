package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handlers) AddTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	insertResult, err := h.collection.InsertOne(context.TODO(), task)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while inserting task in database: %v", err), http.StatusInternalServerError)

		return
	}

	insertedId := insertResult.InsertedID.(bson.ObjectID)
	if err := json.NewEncoder(w).Encode(insertedId); err != nil {
		http.Error(w, fmt.Sprintf("Error writing insertion output: %v", err), http.StatusInternalServerError)
	}
}

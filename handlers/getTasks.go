package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := h.collection.Find(ctx, bson.D{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching tasks from database: %v", err), http.StatusInternalServerError)

		return
	}

	defer cursor.Close(ctx)

	var results []Task
	if err := cursor.All(ctx, &results); err != nil {
		http.Error(w, fmt.Sprintf("Error fetching tasks from database: %v", err), http.StatusInternalServerError)

		return

	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)

		return
	}

}

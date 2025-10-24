package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Kartikeya2710/to-do/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *TaskHandlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	dbCtx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	cursor, err := h.collection.Find(dbCtx, bson.D{})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Fatalf("Error fetching cursor from the collection: %v\n", err)

		return
	}

	defer cursor.Close(dbCtx)

	var results []models.Task
	if err := cursor.All(dbCtx, &results); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Fatalf("Error fetching documents from the cursor: %v\n", err)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Fatalf("Error writing results to response: %v\n", err)

		return
	}

}

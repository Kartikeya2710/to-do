package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	dbCtx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	cursor, err := h.collection.Find(dbCtx, bson.D{})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Fatal("Error fetching cursor from the collection: %v\n", err)

		return
	}

	defer cursor.Close(dbCtx)

	var results []Task
	if err := cursor.All(dbCtx, &results); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Fatal("Error fetching documents from the cursor: %v\n", err)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Fatal("Error writing results to response: %v\n", err)

		return
	}

}

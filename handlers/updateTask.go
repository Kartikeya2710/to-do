package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Kartikeya2710/to-do/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := bson.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		h.logger.Fatalf("Error deriving id from the request: %v\n", err)

		return
	}

	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Bad reqeust", http.StatusBadRequest)
		h.logger.Fatalf("Error deriving task from the request: %v\n", err)

		return
	}

	dbCtx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	if _, err := h.collection.UpdateByID(dbCtx, id, bson.D{{Key: "$set", Value: task}}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Fatalf("Error updating the document in the collection: %v\n", err)

		return
	}

}

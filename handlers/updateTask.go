package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := bson.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Improper format of id specified", http.StatusBadRequest)

		return
	}

	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	dbCtx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	if _, err := h.collection.UpdateByID(dbCtx, id, bson.D{{Key: "$set", Value: task}}); err != nil {
		http.Error(w, fmt.Sprintf("Error updating task in database: %v", err), http.StatusInternalServerError)

		return
	}

}

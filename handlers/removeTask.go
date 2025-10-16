package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handlers) RemoveTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := bson.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)

		return
	}

	filter := bson.D{{Key: "_id", Value: id}}

	dbCtx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	if _, err := h.collection.DeleteOne(dbCtx, filter); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

}

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handlers) RemoveTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := bson.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Improper format of id specified", http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: id}}

	if _, err := h.collection.DeleteOne(ctx, filter); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting task from database: %v", err), http.StatusInternalServerError)

		return
	}

}

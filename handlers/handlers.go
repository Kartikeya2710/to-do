package handlers

import (
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Handlers struct {
	collection *mongo.Collection
	logger     *log.Logger
}

func NewHandlers(collection *mongo.Collection, logger *log.Logger) *Handlers {
	return &Handlers{collection: collection, logger: logger}
}

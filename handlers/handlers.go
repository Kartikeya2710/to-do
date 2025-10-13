package handlers

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Handlers struct {
	collection *mongo.Collection
}

type Task struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string        `bson:"title"`
	Completed bool          `bson:"completed"`
}

func NewHandlers(collection *mongo.Collection) *Handlers {
	return &Handlers{collection: collection}
}

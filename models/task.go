package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Task struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string        `bson:"title"`
	Completed bool          `bson:"completed"`
}

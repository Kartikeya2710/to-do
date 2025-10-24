package handlers

import (
	"log"

	"github.com/Kartikeya2710/to-do/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthHandlers struct {
	collection *mongo.Collection
	jwtManager *utils.JWTManager
	logger     *log.Logger
}

type TaskHandlers struct {
	collection *mongo.Collection
	logger     *log.Logger
}

func NewTaskHandlers(collection *mongo.Collection, logger *log.Logger) *TaskHandlers {
	return &TaskHandlers{collection: collection, logger: logger}
}

func NewAuthHandlers(collection *mongo.Collection, jwtManager *utils.JWTManager, logger *log.Logger) *AuthHandlers {
	return &AuthHandlers{collection: collection, jwtManager: jwtManager, logger: logger}
}

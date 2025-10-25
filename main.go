package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Kartikeya2710/to-do/db"
	"github.com/Kartikeya2710/to-do/handlers"
	"github.com/Kartikeya2710/to-do/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	logger := log.New(os.Stdout, "HTTP: ", log.LstdFlags)
	if err := godotenv.Load("./.env"); err != nil {
		logger.Fatalf("Error loading .env file")
	}

	// DB stuff
	uri, ok := os.LookupEnv("CLUSTER_URI")
	if !ok {
		logger.Fatal("CLUSTER_URI environment variable is not defined")
	}

	dbClient, err := db.NewDBClient(uri, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize DB client: %v", err)
	}
	defer dbClient.Close(context.Background())

	// Tasks DB

	tasksDBName, ok := os.LookupEnv("TASKS_DB_NAME")
	if !ok {
		logger.Fatal("TASKS_DB_NAME environment variable is not defined")
	}

	tasksCollectionName, ok := os.LookupEnv("TASKS_COLLECTION_NAME")
	if !ok {
		logger.Fatal("TASKS_COLLECTION_NAME environment variable is not defined")
	}

	taskCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tasksCollection, err := dbClient.GetMongoDBCollection(taskCtx, tasksDBName, tasksCollectionName)
	if err != nil {
		logger.Fatal("Error fetching MongoDB Collection for tasks")
	}

	// Users DB

	usersDBName, ok := os.LookupEnv("USERS_DB_NAME")
	if !ok {
		logger.Fatal("USERS_DB_NAME environment variable is not defined")
	}

	usersCollectionName, ok := os.LookupEnv("USERS_COLLECTION_NAME")
	if !ok {
		logger.Fatal("USERS_COLLECTION_NAME environment variable is not defined")
	}

	userCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	userCollection, err := dbClient.GetMongoDBCollection(userCtx, usersDBName, usersCollectionName)
	if err != nil {
		logger.Fatal("Error fetching MongoDB Collection for users")
	}

	// JWT Manager

	jwtSecretKey, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		logger.Fatal("JWT_SECRET environment variable is not defined")
	}

	jwtManager := utils.NewJWTManager(jwtSecretKey)

	// Handlers

	authHandlers := handlers.NewAuthHandlers(userCollection, jwtManager, logger)
	taskHandlers := handlers.NewTaskHandlers(tasksCollection, logger)
	router := mux.NewRouter()

	router.HandleFunc("/register", authHandlers.Register).Methods(http.MethodPost)
	router.HandleFunc("/login", authHandlers.Login).Methods(http.MethodPost)

	router.HandleFunc("/tasks", taskHandlers.GetAllTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks", taskHandlers.AddTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks/{id}", taskHandlers.RemoveTask).Methods(http.MethodDelete)
	router.HandleFunc("/tasks/{id}", taskHandlers.UpdateTask).Methods(http.MethodPut)

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Error starting the server: %v", err)
	}
}

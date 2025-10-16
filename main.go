package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Kartikeya2710/to-do/db"
	"github.com/Kartikeya2710/to-do/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	logger := log.New(os.Stdout, "HTTP: ", log.LstdFlags)
	if err := godotenv.Load("./database.env"); err != nil {
		logger.Fatalf("Error loading database.env file")
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		logger.Fatal("DB_NAME environment variable is not defined")
	}

	collectionName, ok := os.LookupEnv("COLLECTION_NAME")
	if !ok {
		logger.Fatal("COLLECTION_NAME environment variable is not defined")
	}

	client, err := db.NewDBClient()
	if err != nil {
		logger.Fatal("Error creating MongoDB client")
	}

	collection, err := db.GetMongoDBCollection(client, dbName, collectionName)
	if err != nil {
		logger.Fatal("Error fetching MongoDB Collection")
	}

	handlers := handlers.NewHandlers(collection, logger)
	router := mux.NewRouter()

	router.HandleFunc("/tasks", handlers.GetAllTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks", handlers.AddTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks/{id}", handlers.RemoveTask).Methods(http.MethodDelete)
	router.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods(http.MethodPut)

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Error starting the server: %v", err)
	}
}

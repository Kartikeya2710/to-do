package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DB struct {
	client *mongo.Client
	logger *log.Logger
}

func NewDBClient(logger *log.Logger) (*DB, error) {
	if _, ok := os.LookupEnv("MONGODB_USER"); !ok {
		return nil, fmt.Errorf("MONGODB_USER environment variable not defined")
	}
	if _, ok := os.LookupEnv("MONGODB_PWD"); !ok {
		return nil, fmt.Errorf("MONGODB_PWD environment variable not defined")
	}

	uri := fmt.Sprintf("mongodb+srv://%s:%s@cluster-0.xzet8ns.mongodb.net/?retryWrites=true&w=majority&appName=Cluster-0", os.Getenv("MONGODB_USER"), os.Getenv("MONGODB_PWD"))

	serverOpts := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverOpts)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB Client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %v", err)
	}

	logger.Println("Connected to MongoDB successfully")

	return &DB{client: client, logger: logger}, nil
}

func (db *DB) GetMongoDBCollection(ctx context.Context, dbName string, collectionName string) (*mongo.Collection, error) {
	database := db.client.Database(dbName)

	collectionMatches, err := database.ListCollectionNames(ctx, map[string]interface{}{"name": collectionName})
	if err != nil {
		db.logger.Fatalf("Error fetching collection with name %s from database %s: %v\n", collectionName, dbName, err)
		return nil, err
	}

	if len(collectionMatches) == 0 {
		db.logger.Printf("Collection %s does not exist. Creating it...\n", collectionName)
		if err := database.CreateCollection(ctx, collectionName); err != nil {
			db.logger.Fatalf("Error creating collection in database %s: %v", dbName, err)
			return nil, err
		}
	}

	db.logger.Printf("Fetched collection %s from database %s successfully\n", collectionName, dbName)

	return database.Collection(collectionName), nil
}

func (db *DB) Close(ctx context.Context) error {
	db.logger.Println("Closing MongoDB connection...")
	
	return db.client.Disconnect(ctx)
}

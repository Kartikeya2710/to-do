package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DBClient struct {
	client *mongo.Client
	logger *log.Logger
}

func NewDBClient(uri string, logger *log.Logger) (*DBClient, error) {
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

	return &DBClient{client: client, logger: logger}, nil
}

func (dbClient *DBClient) GetMongoDBCollection(ctx context.Context, dbName string, collectionName string) (*mongo.Collection, error) {
	database := dbClient.client.Database(dbName)

	collectionMatches, err := database.ListCollectionNames(ctx, map[string]interface{}{"name": collectionName})
	if err != nil {
		dbClient.logger.Fatalf("Error fetching collection with name %s from database %s: %v\n", collectionName, dbName, err)
		return nil, err
	}

	if len(collectionMatches) == 0 {
		dbClient.logger.Printf("Collection %s does not exist. Creating it...\n", collectionName)
		if err := database.CreateCollection(ctx, collectionName); err != nil {
			dbClient.logger.Fatalf("Error creating collection in database %s: %v", dbName, err)
			return nil, err
		}
	}

	dbClient.logger.Printf("Fetched collection %s from database %s successfully\n", collectionName, dbName)

	return database.Collection(collectionName), nil
}

func (dbClient *DBClient) Close(ctx context.Context) error {
	dbClient.logger.Println("Closing MongoDB connection...")

	return dbClient.client.Disconnect(ctx)
}

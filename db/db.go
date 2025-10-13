package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewDBClient() (*mongo.Client, error) {
	if _, ok := os.LookupEnv("MONGODB_USER"); !ok {
		return nil, fmt.Errorf("MONGODB_USER environment variable not defined")
	}
	if _, ok := os.LookupEnv("MONGODB_PWD"); !ok {
		return nil, fmt.Errorf("MONGODB_PWD environment variable not defined")
	}

	uri := fmt.Sprintf("mongodb+srv://%s:%s@cluster-0.xzet8ns.mongodb.net/?retryWrites=true&w=majority&appName=Cluster-0", os.Getenv("MONGODB_USER"), os.Getenv("MONGODB_PWD"))

	serverOpts := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverOpts)

	return mongo.Connect(clientOpts)
}

func GetMongoDBCollection(client *mongo.Client, dbName string, collectionName string) (*mongo.Collection, error) {
	db := client.Database(dbName)

	collectionMatches, err := db.ListCollectionNames(context.TODO(), map[string]interface{}{"name": collectionName})
	if err != nil {
		return nil, err
	}

	if len(collectionMatches) == 0 {
		fmt.Printf("Collection %s does not exist. Creating it...\n", collectionName)
		if err := db.CreateCollection(context.TODO(), collectionName); err != nil {
			return nil, err
		}
	}

	return db.Collection(collectionName), nil
}

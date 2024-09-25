package db

import (
	"go-api/config"

	"go.mongodb.org/mongo-driver/mongo"
)

var Client *mongo.Client = config.ConnectDB()

func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("spark").Collection(collectionName)
}

func ConnectDB() *mongo.Client {
	return Client
}

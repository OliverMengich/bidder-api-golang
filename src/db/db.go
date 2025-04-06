package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func ConnectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/go-mongo-api")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err !=nil{
		log.Fatal(err)
		return nil, err
	}
	log.Println("Connnected to MongoDB")
	return client, nil
}
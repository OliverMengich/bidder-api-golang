package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var (
    Client     *mongo.Client
    DB         *mongo.Database
    AuctionsCol   *mongo.Collection
    ProductsCol *mongo.Collection
    BiddersCol   *mongo.Collection
    BidsCol *mongo.Collection
)

func ConnectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err !=nil{
		log.Fatal(err)
		return nil, err
	}
	log.Println("Connnected to MongoDB")
	Client = client
	DB = client.Database("auction-db") 
	AuctionsCol = DB.Collection("auctions")
	ProductsCol = DB.Collection("products")
	BiddersCol = DB.Collection("bidders")
	BidsCol = DB.Collection("bids")
	return client, nil
}
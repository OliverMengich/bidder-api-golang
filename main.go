package main

import (
	"context"
	"log"
	"time"

	"github.com/OliverMengich/bidder-api-golang/src/db"
	"github.com/OliverMengich/bidder-api-golang/src/handlers"
)

func main()  {
	mongoClient, err := db.ConnectToMongo()
	if err !=nil{
		log.Panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(),15*time.Second)
	defer cancel()
	defer func() {
		if err = mongoClient.Disconnect(ctx); err !=nil {
			panic(err)
		}
	}()
	handlers.CreateRouter()
}
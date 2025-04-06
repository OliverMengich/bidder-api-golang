package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/OliverMengich/bidder-api-golang/src/db"
	"github.com/OliverMengich/bidder-api-golang/src/handlers"
	"github.com/OliverMengich/bidder-api-golang/src/services"
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
	services.New(mongoClient)
	http.ListenAndServe(":3000",nil)
	handlers.CreateRouter()
}
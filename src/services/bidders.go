package services

import (
	"context"
	"log"
	"time"

	"github.com/OliverMengich/bidder-api-golang/src/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bidder struct {
	ID           string               `json:"id,omitempty" bson:"_id,omitempty"`
	BidderNumber float64              `json:"bidder_number,omitempty" bson:"bidder_number,omitempty"`
	Name         string               `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt    time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Updated      time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Auctions     []primitive.ObjectID `json:"auctions" bson:"auctions"`
}

func (b *Bidder) GetAllBidders() ([]Bidder, error) {
	collection := db.BiddersCol
	var bidders []Bidder = []Bidder{}
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var bidder Bidder
		cursor.Decode(&bidder)
		bidders = append(bidders, bidder)
	}
	return bidders, nil
}

func (b *Bidder) RegisterBidder(bidder Bidder) error {
	collection := db.BiddersCol
	_, err := collection.InsertOne(context.TODO(), Bidder{
		Name:         bidder.Name,
		CreatedAt:    time.Now(),
		BidderNumber: bidder.BidderNumber,
		Updated:      time.Now(),
		Auctions:     []primitive.ObjectID{},
	})
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	return nil
}

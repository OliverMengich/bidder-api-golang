package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bidder struct {
	ID           string               `json:"id,omitempty" bson:"_id,omitempty"`
	BidderNumber float64              `json:"bidder_number,omitempty" bson:"bidder_number,omitempty"`
	Name         string               `json:"name,omitempty" bson:"amount,omitempty"`
	CreatedAt    time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Updated      time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Auctions     []primitive.ObjectID `json:"auctions,omitempty" bson:"auctions,omitempty"`
}

func (b *Bidder) GetAllBidders() ([]Bidder, error) {
	collection := returnCollectionPointer("auctions")
	var bidders []Bidder
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
	collection := returnCollectionPointer("bidders")
	_, err := collection.InsertOne(context.TODO(), Bidder{
		Name:         bidder.Name,
		CreatedAt:    time.Now(),
		BidderNumber: bidder.BidderNumber,
		Updated:      time.Now(),
	})
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	return nil
}

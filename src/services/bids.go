package services

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/OliverMengich/bidder-api-golang/src/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bid struct {
	ID        string              `json:"id,omitempty" bson:"_id,omitempty"`
	Amount    float64             `json:"amount,omitempty" bson:"amount,omitempty"`
	CreatedAt time.Time           `json:"created_at,omitempty" bson:"start_time,omitempty"`
	Updated   time.Time           `json:"end_time,omitempty" bson:"end_time,omitempty"`
	AuctionID *primitive.ObjectID `json:"auction_id,omitempty" bson:"auction_id,omitempty"`
	BidderID  *primitive.ObjectID `json:"bidder_id,omitempty" bson:"bidder_id,omitempty"`
}

var bidMutex sync.Mutex

func (b *Bid) GetAllBids() ([]Bid, error) {
	collection := db.BidsCol
	var bids []Bid = []Bid{}
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var bid Bid
		cursor.Decode(&bid)
		bids = append(bids, bid)
	}
	return bids, nil
}
func (b *Bid) PlaceBid(bid Bid) error {
	collection := db.BidsCol
	auctionCollection := db.AuctionsCol
	update := bson.M{"$push": bson.M{"bids": bid.ID}}
	bidMutex.Lock()
	_, err := collection.InsertOne(context.TODO(), Bid{
		Amount:    bid.Amount,
		BidderID:  b.BidderID,
		CreatedAt: time.Now(),
		Updated:   time.Now(),
		AuctionID: bid.AuctionID,
	})
	bidMutex.Unlock()
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	bidMutex.Lock()
	_, err = auctionCollection.UpdateOne(context.Background(), bson.M{"_id": bid.AuctionID}, update)
	bidMutex.Unlock()
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	return nil
}

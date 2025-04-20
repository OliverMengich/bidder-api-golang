package services

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/OliverMengich/bidder-api-golang/src/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Bid struct {
	ID           string              `json:"id,omitempty" bson:"_id,omitempty"`
	Amount       float64             `json:"amount,omitempty" bson:"amount,omitempty"`
	CreatedAt    time.Time           `json:"created_at,omitempty" bson:"start_time,omitempty"`
	Updated      time.Time           `json:"updated_time,omitempty" bson:"updated_time,omitempty"`
	AuctionID    *primitive.ObjectID `json:"auction_id,omitempty" bson:"auction_id,omitempty"`
	Auction      *Auction            `json:"auction"`
	BidderNumber *float64            `json:"bidder_number,omitempty" bson:"bidder_number,omitempty"`
}
type BidInfo struct {
	ID           primitive.ObjectID `json:"id"`
	AuctionID    string             `json:"auction_id"`
	CreatedAt    time.Time          `json:"created_at"`
	Amount       float64            `json:"amount"`
	BidderNumber float64            `json:"bidder_number"`
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
func (b *Bid) PlaceBid(bidAmount float64, bidderNumber *float64, auctionID string) (string, error) {
	collection := db.BidsCol
	auctionCollection := db.AuctionsCol
	mongoID, err := primitive.ObjectIDFromHex(auctionID)
	bidMutex.Lock()
	res, err := collection.InsertOne(context.TODO(), Bid{
		Amount:       bidAmount,
		BidderNumber: bidderNumber,
		CreatedAt:    time.Now(),
		Updated:      time.Now(),
		AuctionID:    &mongoID,
	})

	insertedID := res.InsertedID
	bidID, ok := insertedID.(primitive.ObjectID)
	if !ok {
		log.Println("Failed to convert inserted ID to ObjectID")
		return "", err
	}
	var acc Auction
	updateBids := bson.M{"$push": bson.M{"bids_id": bidID}}
	updateBidders := bson.M{"$addToSet": bson.M{"bidders": bidderNumber}}
	bidMutex.Unlock()
	if err != nil {
		log.Println("Error: ", err)
		return "", err
	}
	bidMutex.Lock()
	_, err = auctionCollection.UpdateOne(context.Background(), bson.M{"_id": mongoID}, updateBids)
	//_, err = auctionCollection.UpdateOne(context.Background(), bson.M{"_id": mongoID}, updateBidders)
	//upsertID := upRes.UpsertedID
	//aucID, ok := upsertID.(primitive.ObjectID)
	err = auctionCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": mongoID}, updateBidders, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&acc)
	bidMutex.Unlock()
	if err != nil {
		log.Println("Error: ", err)
		return "", err
	}
	return auctionID, nil
}

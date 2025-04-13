package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/OliverMengich/bidder-api-golang/src/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionStatus int

const (
	InActive AuctionStatus = iota // 0
	Active                        // 1
	Ended                         // 2
)

type Auction struct {
	ID        string               `json:"id,omitempty" bson:"_id,omitempty"`
	Status    AuctionStatus        `json:"status,omitempty" bson:"status,omitempty"`
	ManagerID string               `json:"manager_id,omitempty" bson:"manager_id,omitempty"`
	StartTime time.Time            `json:"start_time,omitempty" bson:"start_time,omitempty"`
	WinPrice  *float64             `json:"win_price,omitempty" bson:"win_price,omitempty"`
	EndTime   *time.Time           `json:"end_time,omitempty" bson:"end_time,omitempty"`
	WinnerID  *primitive.ObjectID  `json:"winner_id,omitempty" bson:"winner_id,omitempty"`
	ProductID *primitive.ObjectID  `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Bids      []primitive.ObjectID `json:"bids" bson:"bids,omitempty"`
}

func (a *Auction) GetAllAuctions() ([]Auction, error) {
	collection := db.AuctionsCol
	var auctions []Auction = []Auction{}
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Failed fetching", err)
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var auction Auction
		cursor.Decode(&auction)
		auctions = append(auctions, auction)
	}
	return auctions, nil
}
func (a *Auction) CreateAuction(auction Auction) error {
	collection := db.AuctionsCol
	biddersCollection := db.AuctionsCol
	update := bson.M{"$push": bson.M{"bids": auction.ID}}

	_, err := collection.InsertOne(context.TODO(), Auction{
		Status:    Active,
		StartTime: time.Now(),
		ManagerID: auction.ManagerID,
		EndTime:   nil,
		ProductID: auction.ProductID,
	})
	_, err = biddersCollection.UpdateOne(context.Background(), bson.M{"_id": auction.ManagerID}, update)
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	return nil
}
func (a *Auction) GetAuction(auctionID string) (Auction, error) {
	collection := db.AuctionsCol
	mongoID, err := primitive.ObjectIDFromHex(auctionID)
	var auction Auction
	if err != nil {
		log.Fatal(err)
		return Auction{}, err
	}
	err = collection.FindOne(context.Background(), bson.M{"_id": mongoID}).Decode(&auction)
	if err != nil {
		log.Fatal(err)
		return Auction{}, err
	}
	return auction, nil
}
func (a *Auction) EndAuction(auctionID string, entry Auction) (*mongo.UpdateResult, error) {
	collection := db.AuctionsCol
	mongoID, err := primitive.ObjectIDFromHex(auctionID)
	if err != nil {
		return nil, err
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "end_time", Value: time.Now()},
			{Key: "winner_id", Value: entry.WinnerID},
			{Key: "status", Value: Ended},
			{Key: "win_price", Value: entry.WinPrice},
		}},
	}
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": mongoID}, update)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (a *Auction) JoinAuction(auctionID string) (Auction, error) {
	collection := db.AuctionsCol
	mongoID, err := primitive.ObjectIDFromHex(auctionID)
	if err != nil {
		log.Fatal(err)
		return Auction{}, err
	}
	var auction Auction
	err = collection.FindOne(context.Background(), bson.M{"_id": mongoID}).Decode(&auction)
	if err != nil {
		log.Fatal(err)
		return Auction{}, err
	}
	return auction, nil
}

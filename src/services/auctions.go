package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/OliverMengich/bidder-api-golang/src/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuctionStatus int

const (
	InActive AuctionStatus = iota // 0
	Active                        // 1
	Ended                         // 2
)

type Auction struct {
	ID          string               `json:"id,omitempty" bson:"_id,omitempty"`
	Status      AuctionStatus        `json:"status,omitempty" bson:"status,omitempty"`
	ManagerID   float64              `json:"manager_id" bson:"manager_id"`
	StartTime   time.Time            `json:"start_time,omitempty" bson:"start_time,omitempty"`
	WinPrice    *float64             `json:"win_price" bson:"win_price"`
	EndTime     *time.Time           `json:"end_time" bson:"end_time"`
	WinnerID    *primitive.ObjectID  `json:"winner_id" bson:"winner_id"`
	ProductID   *primitive.ObjectID  `json:"product_id" bson:"product_id"`
	BidsID      []primitive.ObjectID `json:"-" bson:"bids_id"`
	Bids        []BidInfor           `json:"bids"`
	Bidders     []float64            `json:"bidders" bson:"bidders"`
	ProductInfo *ProductInfo         `json:"product,omitempty"`
}

type BidInfor struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Amount       float64            `json:"amount" bson:"amount"`
	BidderNumber float64            `json:"bidder_number" bson:"bidder_number"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}
type ProductInfo struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	ReservePrice float64            `json:"reserve_price"`
	ImagesUrl    []string           `json:"images_url"`
}

func (a *Auction) GetAllAuctions() ([]Auction, error) {
	auctionCollection := db.AuctionsCol
	productCollection := db.ProductsCol
	var auctions []Auction = []Auction{}
	cursor, err := auctionCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Failed fetching auctions:", err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var auction Auction
		if err := cursor.Decode(&auction); err != nil {
			log.Println("Decode auction error:", err)
			continue
		}
		// Fetch corresponding product
		var product struct {
			ID           primitive.ObjectID `bson:"_id"`
			Name         string             `bson:"name"`
			ReservePrice float64            `bson:"reserve_price"`
		}
		err := productCollection.FindOne(context.TODO(), bson.M{"_id": auction.ProductID}).Decode(&product)
		if err == nil {
			auction.ProductInfo = &ProductInfo{
				ID:           product.ID,
				Name:         product.Name,
				ReservePrice: product.ReservePrice,
			}
		} else {
			log.Println("Could not find product for auction:", err)
		}
		auctions = append(auctions, auction)
	}
	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}
	return auctions, nil
}
func (a *Auction) CreateAuction(auction Auction) error {
	collection := db.AuctionsCol
	productsCollection := db.ProductsCol
	res, err := collection.InsertOne(context.TODO(), Auction{
		WinPrice:  nil,
		WinnerID:  nil,
		BidsID:    []primitive.ObjectID{},
		Status:    Active,
		StartTime: time.Now(),
		ManagerID: auction.ManagerID,
		EndTime:   nil,
		Bidders:   []float64{},
		ProductID: auction.ProductID,
	})
	insertedID := res.InsertedID 
	auctionID, ok := insertedID.(primitive.ObjectID)
	if !ok {
		log.Println("Failed to convert inserted ID to ObjectID")
		return err
	}
	update := bson.M{"$set": bson.M{"auction_id": auctionID}}

	_, err = productsCollection.UpdateOne(context.Background(), bson.M{"_id": auction.ProductID}, update)
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

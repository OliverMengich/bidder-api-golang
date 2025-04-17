package services

import (
	"context"
	"log"
	"time"

	"github.com/OliverMengich/bidder-api-golang/src/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Product struct {
	ID           string              `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string              `json:"name,omitempty" bson:"name,omitempty"`
	ImagesUrl    []string            `json:"images_url" bson:"images_url"`
	ReservePrice float64             `json:"reserve_price,omitempty" bson:"reserve_price,omitempty"`
	BidderNumber float64             `json:"bidder_number" bson:"bidder_number"`
	CreatedAt    time.Time           `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time           `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	AuctionID    *primitive.ObjectID `json:"auction_id" bson:"auction_id"`
}

var client *mongo.Client

func (p *Product) GetProducts() ([]Product, error) {
	collection := db.ProductsCol

	var products []Product = []Product{}
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var product Product
		cursor.Decode(&product)
		products = append(products, product)
	}
	return products, nil
}
func (p *Product) AddProduct(product Product) error {
	collection := db.ProductsCol
	if product.ImagesUrl == nil {
		product.ImagesUrl = []string{}
	}
	_, err := collection.InsertOne(context.TODO(), Product{
		Name:         product.Name,
		ReservePrice: product.ReservePrice,
		BidderNumber: product.BidderNumber,
		ImagesUrl:    product.ImagesUrl,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		AuctionID:    nil,
	})
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	return nil
}
func (p *Product) GetProductById(id string) (Product, error) {
	collection := db.ProductsCol
	var product Product
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
		return Product{}, err
	}
	err = collection.FindOne(context.Background(), bson.M{"_id": mongoID}).Decode(&product)
	if err != nil {
		log.Fatal(err)
		return Product{}, err
	}
	return product, nil
}
func (p *Product) GetUserProducts(bidderNumber int) ([]Product, error) {
	collection := db.ProductsCol

	var products []Product = []Product{}
	cursor, err := collection.Find(context.TODO(), bson.M{"bidder_number": bidderNumber})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var product Product
		cursor.Decode(&product)
		products = append(products, product)
	}
	return products, nil
}

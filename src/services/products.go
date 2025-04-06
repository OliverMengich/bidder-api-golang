package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Product struct {
	ID           string              `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string              `json:"name,omitempty" bson:"name,omitempty"`
	ReservePrice float64             `json:"reserve_price,omitempty" bson:"reserve_price,omitempty"`
	CreatedAt    time.Time           `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time           `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	AuctionID    *primitive.ObjectID `json:"auction_id,omitempty" bson:"auction_id,omitempty"`
}

var client *mongo.Client

func New(mongo *mongo.Client) Product {
	client = mongo
	return Product{}
}
func returnCollectionPointer(collection string) *mongo.Collection {
	return client.Database("auction_db").Collection(collection)
}
func (p *Product) GetProducts() ([]Product, error) {
	collection := returnCollectionPointer("products")
	var products []Product
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
	collection := returnCollectionPointer("products")
	_, err := collection.InsertOne(context.TODO(), Product{
		Name:         product.Name,
		ReservePrice: product.ReservePrice,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	})
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	return nil
}
func (p *Product) GetProductById(id string) (Product, error) {
	collection := returnCollectionPointer("products")
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

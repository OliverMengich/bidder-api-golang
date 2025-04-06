package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OliverMengich/bidder-api-golang/src/services"
)

var product services.Product

func healthCheck(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Msg:  "Health Check",
		Code: 200,
	}
	respondWithJSON(w, 200, res)
}
func getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := product.GetProducts()
	if err != nil {
		responseWithError(w, 400, err.Error())
	}
	respondWithJSON(w, 200, products)
}
func createProduct(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		responseWithError(w, 400, "Error adding Product")
	}

	err = product.AddProduct(product)
	if err != nil {
		responseWithError(w, 400, "Error adding Product")
	}
	res := Response{
		Msg:  "Successfully Added Product to auction",
		Code: 201,
	}
	respondWithJSON(w, res.Code, res)
}

func getProductById(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("productID")
	fmt.Println(productID)
	product, err := product.GetProductById(productID)
	if err != nil {
		responseWithError(w, 400, "Could not find Product")
	}
	respondWithJSON(w, 200, product)
}

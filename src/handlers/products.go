package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/OliverMengich/bidder-api-golang/src/services"
)

var product services.Product

func getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := product.GetProducts()
	if err != nil {
		responseWithError(w, 400, err.Error())
		return
	}
	respondWithJSON(w, 200, products)
}
func createProduct(w http.ResponseWriter, r *http.Request) {
	uploadDir := "src/handlers/uploads"
	r.ParseMultipartForm(10 << 20)
	productJSON := r.FormValue("product")

	err := json.Unmarshal([]byte(productJSON), &product)
	if err != nil {
		fmt.Println("Error decode", err)
		responseWithError(w, 400, "Error adding Product")
		return
	}
	// Handle multiple images
	formdata := r.MultipartForm
	files := formdata.File["images"] // match "images" with app key
	images := []string{}
	for _, handler := range files {
		fmt.Println("Image file:", handler.Filename)
		src, err := handler.Open()
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}
		defer src.Close()

		// Make sure upload dir exists
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			_ = os.MkdirAll(uploadDir, os.ModePerm)
		}

		dstPath := filepath.Join(uploadDir, handler.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			fmt.Println("OS create error: ", err)
			continue
		}
		defer dst.Close()
		images = append(images, handler.Filename)
		_, err = io.Copy(dst, src)
		if err != nil {
			fmt.Println("Error copying:", err)
		}
	}

	err = product.AddProduct(services.Product{
		Name:         product.Name,
		ImagesUrl:    images,
		ReservePrice: product.ReservePrice,
		BidderNumber: product.BidderNumber,
	})
	images = []string{}
	if err != nil {
		fmt.Println("adding error:", err)
		responseWithError(w, 400, "Error adding Product")
		return
	}
	res := Response{
		message: "Successfully Added Product to auction",
		code:    201,
	}
	respondWithJSON(w, res.code, res)
}

func getProductById(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("productID")
	fmt.Println(productID)
	product, err := product.GetProductById(productID)
	if err != nil {
		responseWithError(w, 400, "Could not find Product")
		return
	}
	respondWithJSON(w, 200, product)
}
func getUserProducts(w http.ResponseWriter, r *http.Request) {
	bidderNumber, err := strconv.Atoi(r.PathValue("bidderNumber"))
	if err != nil {
		responseWithError(w, 400, "Could not find Product")
	}
	fmt.Println("Bidder number: ", bidderNumber)
	products, err := product.GetUserProducts(bidderNumber)
	if err != nil {
		responseWithError(w, 400, err.Error())
		return
	}
	respondWithJSON(w, 200, products)
}

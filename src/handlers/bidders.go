package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/OliverMengich/bidder-api-golang/src/services"
)

var bidder services.Bidder

func registerBidder(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&bidder)
	if err != nil {
		log.Println(err)
		responseWithError(w, 400, "Error adding Product")
		return
	}

	err = bidder.RegisterBidder(bidder)
	if err != nil {
		responseWithError(w, 400, "Error registering Bidder ")
		return
	}
	res := Response{
		message: "Successfully Registered new member",
		code:    201,
	}
	respondWithJSON(w, res.code, res)
}

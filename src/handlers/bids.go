package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/OliverMengich/bidder-api-golang/src/services"
	"golang.org/x/net/websocket"
)

var bid services.Bid

func getAllBids(w http.ResponseWriter, r *http.Request) {

}
func getBidsByAuctionID(w http.ResponseWriter, r *http.Request) {

}
func placeBid(ws *websocket.Conn) {
	body := ws.Request().Body
	err := json.NewDecoder(body).Decode(&bid)
	if err != nil {
		log.Fatal(err)
	}
}

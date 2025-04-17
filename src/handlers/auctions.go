package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/OliverMengich/bidder-api-golang/src/services"
	"golang.org/x/net/websocket"
)

var auction services.Auction

func getAuctions(w http.ResponseWriter, r *http.Request) {
	auctions, err := auction.GetAllAuctions()
	if err != nil {
		responseWithError(w, 400, err.Error())
		return
	}
	respondWithJSON(w, 200, auctions)
}
func createAuction(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&auction)
	if err != nil {
		responseWithError(w, 400, "Error adding auction")
		return
	}

	err = auction.CreateAuction(auction)
	if err != nil {
		responseWithError(w, 400, "Error adding auction")
		return
	}
	res := Response{
		message: "Successfully Added auction to auction",
		code:    201,
	}
	respondWithJSON(w, res.code, res)
}
func getAuction(w http.ResponseWriter, r *http.Request) {
	auctionID := r.PathValue("auctionID")
	auction, err := auction.GetAuction(auctionID)
	if err != nil {
		responseWithError(w, 404, "Auction Not found")
		return
	}
	respondWithJSON(w, 200, auction)
}

func readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Read error:", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))
		ws.Write([]byte("Thank you for the message !!"))
		// s.broadcast(msg)
	}
}

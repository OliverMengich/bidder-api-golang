package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/OliverMengich/bidder-api-golang/src/services"
	"golang.org/x/net/websocket"
)

type WebsocketMessage struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

type GenericStructCustom interface {
	string | services.BidInfo
}
type WebsocketResponse[T GenericStructCustom] struct {
	Status  string `json:"status"`
	Type    string `json:"type"`
	Payload T      `json:"payload"`
}
type Response struct {
	message string
	code    int
}
type APIServer struct {
	addr  string
	conns map[*websocket.Conn]bool
	mu    sync.Mutex
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr:  addr,
		conns: make(map[*websocket.Conn]bool),
	}
}
func (s *APIServer) Run() error {
	fs := http.FileServer(http.Dir("./src/handlers/uploads"))
	router := http.NewServeMux()
	router.Handle("/images/", http.StripPrefix("/images/", fs)) // serve the images
	router.Handle("/auction-socket", websocket.Handler(s.handleWS))
	router.HandleFunc("GET /products", getProducts)
	router.HandleFunc("GET /products/owner/{bidderNumber}", getUserProducts)
	router.HandleFunc("POST /products", createProduct)
	router.HandleFunc("GET /products/{productID}", getProductById)
	router.HandleFunc("GET /auctions", getAuctions)
	router.HandleFunc("POST /auctions", createAuction)
	router.HandleFunc("POST /bidders", registerBidder)
	router.HandleFunc("GET /auctions/{auctionID}", getAuction)

	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}
	log.Printf("Server has started %s", s.addr)
	return server.ListenAndServe()
}
func (s *APIServer) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client: ", ws.RemoteAddr())
	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *APIServer) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	defer func() {
		a := &WebsocketResponse[string]{Status: "success", Type: "client_disconnect", Payload: "Client disconnected: "+ws.RemoteAddr().String()}
		out, err := json.Marshal(a)
		if err != nil {
			panic(err)
		}
		s.broadcast([]byte(out))
		log.Fatal(err)
		log.Println("Client disconnected")
		ws.Close()
		return
	}()
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

		var message WebsocketMessage
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("JSON Unmarshalling error:", err)
			continue
		}
		switch message.Type {
		case "end_auction":
			//expects winner id, winning price
			auctionID := message.Payload["auction_id"].(string)
			winnerID := message.Payload["winner_id"].(string)
			winnerPrice := message.Payload["win_price"].(float64)
			_, err := auction.EndAuction(auctionID, winnerID, winnerPrice)
			if err != nil {
				a := &WebsocketResponse[string]{Status: "success", Type: "end_auction", Payload: "Could not End Auction " + err.Error()}
				out, err := json.Marshal(a)
				if err != nil {
					panic(err)
				}
				s.broadcast([]byte(out))
				log.Fatal(err)
				return
			}
			a := &WebsocketResponse[string]{Type: "end_auction", Status: "success", Payload: "Auction ended successfully"}
			out, err := json.Marshal(a)
			if err != nil {
				panic(err)
			}
			s.broadcast([]byte(out))
			log.Fatal(err)
			return
		case "join_auction":
			auctionID := message.Payload["auction_id"].(string)
			auction, err := auction.GetAuction(auctionID)
			if err != nil {
				a := &WebsocketResponse[string]{Type: "join_auction", Status: "failure", Payload: "Could not find Auction " + err.Error()}
				out, err := json.Marshal(a)
				if err != nil {
					panic(err)
				}
				s.broadcast([]byte(out))
				log.Fatal(err)
				return
			}
			if auction.Status == 2 {
				a := &WebsocketResponse[string]{Type: "join_auction", Status: "failure", Payload: "Auction has ended"}
				out, err := json.Marshal(a)
				if err != nil {
					panic(err)
				}
				s.broadcast([]byte(out))
				log.Fatal("Auction has ended")
				return
			}
			a := &WebsocketResponse[string]{Type: "join_auction", Status: "join_auction", Payload: "Successfully joined auction"}
			out, err := json.Marshal(a)
			if err != nil {
				panic(err)
			}
			s.broadcast([]byte(out))
			fmt.Println("Joining auction: ", message.Payload)
		case "place_bid":
			// expects auctionID, bid amount, bidder ID
			fmt.Println("auction id: ", message.Payload["auction_id"])
			fmt.Println("bid amount: ", message.Payload["bid_amount"])
			fmt.Println("bidder number: ", message.Payload["bidder_number"])
			auctionID := message.Payload["auction_id"]
			bidAmount := message.Payload["bid_amount"]
			bidderNumber := message.Payload["bidder_number"]

			bidAmountFloat, ok := bidAmount.(float64)
			if !ok {
				a := &WebsocketResponse[string]{Type: "place_bid", Status: "failure", Payload: "Invalid bid amount"}
				out, err := json.Marshal(a)
				if err != nil {
					panic(err)
				}
				s.broadcast([]byte(out))
				return
			}
			bidderNumberFloat, ok := bidderNumber.(float64)
			if !ok {
				a := &WebsocketResponse[string]{Type: "place_bid", Status: "failure", Payload: "Invalid bidder number"}
				out, err := json.Marshal(a)
				if err != nil {
					panic(err)
				}
				s.broadcast([]byte(out))
				return
			}
			bidInfo, err := bid.PlaceBid(bidAmountFloat, &bidderNumberFloat, auctionID.(string))
			if err != nil {
				fmt.Println("Error placing bid: ", err)
				a := &WebsocketResponse[string]{Type: "place_bid", Status: "success", Payload: "Could not place bid"}
				out, err := json.Marshal(a)
				if err != nil {
					panic(err)
				}
				s.broadcast([]byte(out))
				return
			}
			a := &WebsocketResponse[string]{Type: "place_bid", Status: "success", Payload: bidInfo}
			out, err := json.Marshal(a)
			if err != nil {
				panic(err)
			}
			s.broadcast([]byte(out))
			fmt.Println("placing a bid:", message.Payload)
		default:
			fmt.Println("Unknown message type:", message.Type)
		}
		fmt.Println(string(msg))
		//ws.Write([]byte("Thank you for the message !!"))
		//s.broadcast(msg)
	}
}
func (s *APIServer) broadcast(b []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("Writing error:", err)
				ws.Close()
				s.mu.Lock()
				delete(s.conns, ws)
				s.mu.Unlock()
			}
		}(conn)
	}
}
func CreateRouter() error {
	server := NewAPIServer(":3000")
	return server.Run()
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
func responseWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("5xx error:", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

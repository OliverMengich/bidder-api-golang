package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Response struct {
	Msg  string
	Code int
}
type APIServer struct {
	addr  string
	conns map[*websocket.Conn]bool
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr:  addr,
		conns: make(map[*websocket.Conn]bool),
	}
}
func (s *APIServer) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /products", getProducts)
	router.HandleFunc("POST /products", createProduct)
	router.HandleFunc("GET /products/{productID}", getProductById)
	router.HandleFunc("GET /auctions", getAuctions)
	router.HandleFunc("GET /auctions/{auctionID}", getAuction)
	router.HandleFunc("POST /auctions/end", func(w http.ResponseWriter, r *http.Request) {
		websocket.Handler(endAuction).ServeHTTP(w, r)
	})
	router.HandleFunc("GET /auctions/join/{auctionID}", func(w http.ResponseWriter, r *http.Request) {
		websocket.Handler(joinAuction).ServeHTTP(w, r)
	})
	router.HandleFunc("POST /bids/place/{auctionID}", func(w http.ResponseWriter, r *http.Request) {
		websocket.Handler(placeBid).ServeHTTP(w, r)
	})

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
func (s *APIServer) handleJoinAuction(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client: ", ws.RemoteAddr())
	s.conns[ws] = true
	s.readLoop(ws)
}
func (s *APIServer) readLoop(ws *websocket.Conn) {
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
		s.broadcast(msg)
	}
}
func (s *APIServer) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("Writing error: ", err)
			}
		}(ws)
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

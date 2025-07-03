package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for simplicity; in production, you should implement proper origin checks.
		return true
	},
}

func HandlerWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Handle the WebSocket connection here.
	// For example, you can read messages from the client and respond.
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break // Exit on error
		}
		log.Printf("Received : %s", string(msg))

	}
}

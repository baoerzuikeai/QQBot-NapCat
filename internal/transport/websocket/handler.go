package websocket

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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
		setStatusCode(w, 10)
	}
}

func setStatusCode(w http.ResponseWriter, code int) {
	url := "http://192.168.2.190:3000/send_private_msg"
	method := "POST"

	payload := strings.NewReader(`{
    "user_id": "2634174807",
    "message": [
        {
            "type": "text",
            "data": {
                "text": "napcat"
            }
        }
    ]
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

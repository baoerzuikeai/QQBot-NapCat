package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/baoerzuikeai/QQBot-NapCat/internal/aiclient"
	"github.com/baoerzuikeai/QQBot-NapCat/internal/onebot"
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

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break // Exit on error
		}
		var baseEvent onebot.BaseEvent
		err = json.Unmarshal(msg, &baseEvent)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue // Skip this message if unmarshalling fails

		}
		switch baseEvent.PostType {
		case "message":
			var messageEvent onebot.MessageEvent
			err = json.Unmarshal(msg, &messageEvent)
			if err != nil {
				log.Printf("Error unmarshalling message event: %v", err)
				continue // Skip this message if unmarshalling fails
			}
			log.Printf("Received message event: %+v", messageEvent)
			switch messageEvent.MessageType {
			case "private":
				var privateMessageEvent onebot.PrivateMessageEvent
				err = json.Unmarshal(msg, &privateMessageEvent)
				if err != nil {
					log.Printf("Error unmarshalling private message event: %v", err)
					continue // Skip this message if unmarshalling fails
				}
				sendPrivateMsg(strconv.FormatInt(privateMessageEvent.UserID, 10), privateMessageEvent.RawMessage)

			}

		default:
			log.Println("Unknown post type:", baseEvent.PostType)
		}
	}
}

func sendPrivateMsg(userid, msg string) {
	url := "http://localhost:3000/send_private_msg"
	method := "POST"

	deepSeek := aiclient.NewDeepSeekClient("sk-78305562a4c942748da616a3d9a58ad7", "https://api.deepseek.com/chat/completions")
	//读取txt文本
	contentBytes, _ := os.ReadFile("system.txt")
	history := []aiclient.Message{
		{
			Role:    "system",
			Content: string(contentBytes),
		},
	}
	his, _ := deepSeek.GetResponse(context.Background(), history, msg)

	payload := strings.NewReader(`{
    "user_id": ` + `"` + userid + `"` + `,
    "message": [
        {
            "type": "text",
            "data": {
                "text": ` + `"` + his[0].Message.Content + `"` + `
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
	log.Println(string(body))
}

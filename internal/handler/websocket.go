package websocket

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/baoerzuikeai/QQBot-NapCat/internal/aiclient"
	"github.com/baoerzuikeai/QQBot-NapCat/internal/onebot"
	"github.com/baoerzuikeai/QQBot-NapCat/internal/repository/sqlite"
	"github.com/baoerzuikeai/QQBot-NapCat/internal/service"
	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	AiChatService service.AIChatServiceInterface
}

func NewWebsocketHandler(aiChatService service.AIChatServiceInterface) *WebsocketHandler {
	return &WebsocketHandler{
		AiChatService: aiChatService,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for simplicity; in production, you should implement proper origin checks.
		return true
	},
}

func (wh *WebsocketHandler) HandlerWebsocket(w http.ResponseWriter, r *http.Request) {
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
		wh.despatchOneBotEvent(msg)
	}
}

func (wh *WebsocketHandler) despatchOneBotEvent(eventByte []byte) {
	var baseEvent onebot.BaseEvent
	err := json.Unmarshal(eventByte, &baseEvent)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return // Skip this message if unmarshalling fails

	}
	switch baseEvent.PostType {
	case "message":
		var messageEvent onebot.MessageEvent
		err = json.Unmarshal(eventByte, &messageEvent)
		if err != nil {
			log.Printf("Error unmarshalling message event: %v", err)
			return // Skip this message if unmarshalling fails
		}
		log.Printf("Received message event: %+v", messageEvent)
		switch messageEvent.MessageType {
		case "private":
			var privateMessageEvent onebot.PrivateMessageEvent
			err = json.Unmarshal(eventByte, &privateMessageEvent)
			if err != nil {
				log.Printf("Error unmarshalling private message event: %v", err)
				return // Skip this message if unmarshalling fails
			}
			wh.SendPrivateMsg(privateMessageEvent.UserID, privateMessageEvent.RawMessage)

		}

	default:
		log.Println("Unknown post type:", baseEvent.PostType)
	}
}

func (wh *WebsocketHandler) SendPrivateMsg(userid int64, msg string) {
	url := "http://localhost:3000/send_private_msg"
	method := "POST"
	var sqhis []sqlite.AIHistory
	//读取txt文本
	sqhis, err := wh.AiChatService.GetAIHistoryByUserID(userid)
	if err != nil {
		log.Printf("Error getting AI history: %v", err)
		return
	}
	var history []aiclient.Message
	if len(sqhis) == 0 {
		for _, his := range sqhis {
			history = append(history, aiclient.Message{
				Role:    his.Role,
				Content: his.Content,
			})
		}
	}
	hisc, err := wh.AiChatService.GetResponse(history, msg)
	wh.AiChatService.SaveAIhistory(userid, hisc[0].Message.Content, hisc[0].Message.Role, "session_id")
	wh.AiChatService.SaveAIhistory(userid, msg, "user", "session_id")
	payload := strings.NewReader(`{
    "user_id": ` + `"` + strconv.FormatInt(userid, 10) + `"` + `,
    "message": [
        {
            "type": "text",
            "data": {
                "text": ` + `"` + hisc[0].Message.Content + `"` + `
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

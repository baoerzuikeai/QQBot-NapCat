package main

import (
	"net/http"

	"github.com/baoerzuikeai/QQBot-NapCat/internal/transport/websocket"
)

func main() {
	http.HandleFunc("/ws/napcat", websocket.HandlerWebsocket)
	http.ListenAndServe("0.0.0.0:8082", nil)
}

package main

import (
	"net/http"

	websocket "github.com/baoerzuikeai/QQBot-NapCat/internal/handler"
)

func main() {
	http.HandleFunc("/ws/napcat", websocket.HandlerWebsocket)
	http.ListenAndServe("0.0.0.0:8082", nil)
}

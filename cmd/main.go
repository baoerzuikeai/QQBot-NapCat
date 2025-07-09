package main

import (
	"net/http"

	"github.com/baoerzuikeai/QQBot-NapCat/internal/aiclient"
	websocket "github.com/baoerzuikeai/QQBot-NapCat/internal/handler"
	"github.com/baoerzuikeai/QQBot-NapCat/internal/repository/sqlite"
	"github.com/baoerzuikeai/QQBot-NapCat/internal/service"
)

const dbFileName = "ai_history.db"

func main() {
	aihistoryrepo, _ := sqlite.NewSqliteAIHistoryRepository(dbFileName)
	deepSeek := aiclient.NewDeepSeekClient("sk-78305562a4c942748da616a3d9a58ad7", "https://api.deepseek.com/chat/completions")
	aichatsrv := service.NewAIChatService(aihistoryrepo, deepSeek)
	wh := websocket.NewWebsocketHandler(aichatsrv)
	http.HandleFunc("/ws/napcat", wh.HandlerWebsocket)
	http.ListenAndServe("0.0.0.0:8082", nil)
}

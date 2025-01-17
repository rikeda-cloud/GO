package handlers

import (
	"github.com/gorilla/websocket"
)

type WebSocketHandler interface {
	ReadFromWebSocket(ws *websocket.Conn) error
	WriteToWebSocket(ws *websocket.Conn) error
}

type WebSocketBaseHandler struct {
	Upgrader websocket.Upgrader
}

func NewWebSocketBaseHandler() *WebSocketBaseHandler {
	return &WebSocketBaseHandler{
		Upgrader: websocket.Upgrader{},
	}
}

const (
	NORMAL = "NORMAL"
	DELETE = "DELETE"
	FINISH = "FINISH"
	NEXT   = "NEXT"
	PREV   = "PREV"
	MOD    = "MOD"
)

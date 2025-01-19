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

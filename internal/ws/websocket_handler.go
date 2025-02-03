package ws

import (
	"github.com/gorilla/websocket"
)

type WebSocketBaseHandler struct {
	Upgrader websocket.Upgrader
}

func NewWebSocketBaseHandler() *WebSocketBaseHandler {
	return &WebSocketBaseHandler{
		Upgrader: websocket.Upgrader{},
	}
}

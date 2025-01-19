package handlers

import (
	"GO/internal/point"
	"encoding/json"
	"github.com/gorilla/websocket"
)

const (
	NORMAL = "NORMAL"
	DELETE = "DELETE"
	FINISH = "FINISH"
	NEXT   = "NEXT"
	PREV   = "PREV"
	MOD    = "MOD"
)

type ImageMarkData struct {
	FileName string      `json:"file_name"`
	Point    point.Point `json:"point"`
	Control  string      `json:"control"`
	Tags     string      `json:"tags"`
}

type RemainImageCountData struct {
	CurrentCount  int `json:"current_count"`
	PreviousCount int `json:"previous_count"`
}

func SendCarData(ws *websocket.Conn, fileName string, point point.Point, control string) error {
	sendData := ImageMarkData{
		FileName: fileName,
		Point:    point,
		Control:  control,
		Tags:     "",
	}
	data, err := json.Marshal(sendData)
	if err != nil {
		return err
	}
	return ws.WriteMessage(websocket.TextMessage, data)
}

func SendRemainImageCountData(ws *websocket.Conn, curentCount, previousCount int) error {
	sendData := RemainImageCountData{
		CurrentCount:  curentCount,
		PreviousCount: previousCount,
	}
	data, err := json.Marshal(sendData)
	if err != nil {
		return err
	}
	return ws.WriteMessage(websocket.TextMessage, data)
}

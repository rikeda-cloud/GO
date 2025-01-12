package handlers

import (
	"GO/internal/db"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type RemainImageCountHandler struct {
	WebSocketBaseHandler
	Count int
}

type RemainImageCountData struct {
	CurrentCount  int `json:"current_count"`
	PreviousCount int `json:"previous_count"`
}

func NewRemainImageCountHandler() *RemainImageCountHandler {
	count, err := carDataDB.SelectRemainImageCount()
	if err != nil {
		panic("Error: Get Remain Image Count")
	}
	return &RemainImageCountHandler{
		WebSocketBaseHandler: *NewWebSocketBaseHandler(),
		Count:                count,
	}
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

func (wsh *RemainImageCountHandler) RemainImageCountHandler(c echo.Context) error {
	conn, err := wsh.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	SendRemainImageCountData(conn, wsh.Count, wsh.Count)

	for {
		// Write
		if err := wsh.WriteToWebSocket(conn); err != nil {
			c.Logger().Info(err)
		}
	}
}

func (wsh *RemainImageCountHandler) WriteToWebSocket(ws *websocket.Conn) error {
	count, err := carDataDB.SelectRemainImageCount()
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 100)

	// 残り画像枚数に変化がなければ何もしない
	if wsh.Count == count {
		return nil
	}
	errOrNil := SendRemainImageCountData(ws, count, wsh.Count)
	wsh.Count = count
	return errOrNil
}
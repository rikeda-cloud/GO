package handlers

import (
	"GO/internal/db"
	"time"

	"github.com/gorilla/websocket"
)

type PredictedRemainImageCountHandler struct {
	RemainImageCountHandler
	Count int
}

func NewPredictedRemainImageCountHandler() *PredictedRemainImageCountHandler {
	count, err := carDataDB.SelectPredictedRemainImageCount()
	if err != nil {
		panic("Error: Get Predicted Remain Image Count")
	}
	return &PredictedRemainImageCountHandler{
		RemainImageCountHandler: *NewRemainImageCountHandler(),
		Count:                   count,
	}
}

func (wsh *PredictedRemainImageCountHandler) WriteToWebSocket(ws *websocket.Conn) error {
	count, err := carDataDB.SelectPredictedRemainImageCount()
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 300)

	// 残り画像枚数に変化がなければ何もしない
	if wsh.Count == count {
		return nil
	}
	errOrNil := SendRemainImageCountData(ws, count, wsh.Count)
	wsh.Count = count
	return errOrNil
}

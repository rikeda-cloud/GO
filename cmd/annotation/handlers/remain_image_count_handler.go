package handlers

import (
	"GO/internal/db"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type RemainImageCountHandler struct {
	WebSocketBaseHandler
	Count int
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

func (wsh *RemainImageCountHandler) HandleRemainImageCount(c echo.Context) error {
	conn, err := wsh.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	SendRemainImageCountData(conn, wsh.Count, wsh.Count)

	// ゴルーチンで非ブロッキングな読み取り処理を実行
	done := make(chan struct{}) // ゴルーチン終了を通知するチャネル
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("[Remain Read]: ", err)
				close(done)
				break
			}
		}
	}()

	for {
		// Write
		if err := wsh.WriteToWebSocket(conn); err != nil {
			fmt.Println("[Remain Write]: ", err)
			return nil
		}

		// ゴルーチンが終了したかを待つ
		select {
		case <-done:
			// ゴルーチンが終了
			return nil
		default:
			// ゴルーチンが終了していない場合は何もしない
		}
	}
}

func (wsh *RemainImageCountHandler) WriteToWebSocket(ws *websocket.Conn) error {
	count, err := carDataDB.SelectRemainImageCount()
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

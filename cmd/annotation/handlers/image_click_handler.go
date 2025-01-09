package handlers

import (
	"GO/internal/config"
	"GO/internal/db"
	"GO/internal/point"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ImageClickHandler struct {
	WebSocketBaseHandler
	PrevDataId       int64
	BasePoint        point.Point
	MaxDistancePoint point.Point
}

func NewImageClickHandler() *ImageClickHandler {
	cfg := config.GetConfig()
	return &ImageClickHandler{
		WebSocketBaseHandler: *NewWebSocketBaseHandler(),
		PrevDataId:           0,
		BasePoint:            point.Point{X: float64(cfg.Camera.Width / 2), Y: float64(cfg.Camera.Height)},
		MaxDistancePoint:     point.Point{X: 0, Y: 0},
	}
}

const (
	NORMAL = "NORMAL"
	DELETE = "DELETE"
	FINISH = "FINISH"
)

type ImageMarkData struct {
	FileName string      `json:"file_name"`
	Point    point.Point `json:"point"`
	Control  string      `json:"control"`
}

func (wsh *ImageClickHandler) ImageClickHandler(c echo.Context) error {
	conn, err := wsh.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		// Write
		if err := wsh.WriteToWebSocket(conn); err != nil {
			c.Logger().Info(err)
		}

		// Read
		if err := wsh.ReadFromWebSocket(conn); err != nil {
			c.Logger().Info(err)
		}
	}
}

func SendCarData(ws *websocket.Conn, fileName string, point point.Point, control string) error {
	sendData := ImageMarkData{
		FileName: fileName,
		Point:    point,
		Control:  control,
	}
	data, err := json.Marshal(sendData)
	if err != nil {
		return err
	}
	return ws.WriteMessage(websocket.TextMessage, data)
}

func (wsh *ImageClickHandler) WriteToWebSocket(ws *websocket.Conn) error {
	carData, err := carDataDB.SelectNoMarkedCarData(wsh.PrevDataId)

	// 全てのデータがアノテーション済み
	if err == sql.ErrNoRows {
		return SendCarData(ws, "", point.Point{X: 0, Y: 0}, FINISH)
	}
	if err != nil {
		return err
	}
	// 取得したIDを保存、次回、保存したIDより１つ大きいIDをSELECT
	wsh.PrevDataId = int64(carData.ID)

	magnitude := carData.CarSpeed
	angle := -carData.CarSteering
	actPoint := point.ReverseCalculate(wsh.BasePoint, wsh.MaxDistancePoint, magnitude, angle)

	return SendCarData(ws, carData.FileName, actPoint, NORMAL)
}

func (wsh *ImageClickHandler) ReadFromWebSocket(ws *websocket.Conn) error {
	_, msg, err := ws.ReadMessage()
	if err != nil {
		return err
	}

	var data ImageMarkData
	err = json.Unmarshal(msg, &data)
	if data.Control == DELETE {
		return carDataDB.DeleteCarData(data.FileName)
	}
	if err != nil {
		return err
	}

	clickPoint := data.Point
	magnitude := point.CalcNormalizedMagnitude(wsh.BasePoint, clickPoint, wsh.MaxDistancePoint)
	angle := -(point.CalcAngle(wsh.BasePoint, clickPoint))
	reverse := point.ReverseCalculate(wsh.BasePoint, wsh.MaxDistancePoint, magnitude, angle)

	fmt.Println("magnitude: ", magnitude)
	fmt.Println("Angle   : ", int(angle))
	fmt.Println("reverse : ", reverse)

	return carDataDB.UpdateCarData(data.FileName, magnitude, angle)
}

package handlers

import (
	"GO/internal/db"
	"GO/internal/point"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type PredictedDataHandler struct {
	WebSocketBaseHandler
	CoordinateRange
	PrevDataId int64
}

func NewPredictedDataHandler() *PredictedDataHandler {
	return &PredictedDataHandler{
		WebSocketBaseHandler: *NewWebSocketBaseHandler(),
		CoordinateRange:      *NewCoordinateRange(),
		PrevDataId:           0,
	}
}

func (wsh *PredictedDataHandler) HandlePredictedData(c echo.Context) error {
	conn, err := wsh.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		// Write
		if err := wsh.WriteToWebSocket(conn); err != nil {
			fmt.Println("[Write]: ", err)
			return nil
		}

		// Read
		if err := wsh.ReadFromWebSocket(conn); err != nil {
			fmt.Println("[Read]: ", err)
			return nil
		}
	}
}

func (wsh *PredictedDataHandler) WriteToWebSocket(ws *websocket.Conn) error {
	carData, err := carDataDB.SelectPredictedNoMarkedCarData(wsh.PrevDataId)

	if err == sql.ErrNoRows {
		wsh.PrevDataId = 0
		carData, err = carDataDB.SelectPredictedNoMarkedCarData(wsh.PrevDataId)
		if err == sql.ErrNoRows {
			return SendCarData(ws, "", point.Point{X: 0, Y: 0}, FINISH)
		}
	}
	if err != nil {
		return err
	}
	wsh.PrevDataId = int64(carData.ID)

	magnitude := carData.CarSpeed
	angle := carData.CarSteering
	predictedPoint := point.ReverseCalculate(wsh.BasePoint, wsh.MaxDistancePoint, magnitude, angle)
	return SendCarData(ws, carData.FileName, predictedPoint, NORMAL)
}

func (wsh *PredictedDataHandler) ReadFromWebSocket(ws *websocket.Conn) error {
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
	tags := data.Tags

	fmt.Println("magnitude: ", magnitude)
	fmt.Println("Angle    : ", int(angle))
	fmt.Println("reverse  : ", reverse)
	fmt.Println("tags     : ", tags)

	return carDataDB.UpdateCarData(data.FileName, magnitude, angle, tags)
}

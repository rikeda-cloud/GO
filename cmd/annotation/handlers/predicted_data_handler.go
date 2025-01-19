package handlers

import (
	"GO/internal/db"
	"GO/internal/point"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type PredictedDataHandler struct {
	ImageClickHandler
}

func NewPredictedDataHandler() *PredictedDataHandler {
	return &PredictedDataHandler{
		ImageClickHandler: *NewImageClickHandler(),
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

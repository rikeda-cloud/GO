package handlers

import (
	"GO/internal/db"
	"GO/internal/point"
	"GO/internal/ws"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type AnnotatedDataCheckHandler struct {
	ws.WebSocketBaseHandler
	CoordinateRange
}

func NewAnnotatedDataCheckHandler() *AnnotatedDataCheckHandler {
	return &AnnotatedDataCheckHandler{
		WebSocketBaseHandler: *ws.NewWebSocketBaseHandler(),
		CoordinateRange:      *NewCoordinateRange(),
	}
}

type AnnotateData struct {
	FileName       string      `json:"file_name"`
	ActPoint       point.Point `json:"act_point"`
	AnnotatedPoint point.Point `json:"annotated_point"`
	Control        string      `json:"control"`
	Tags           string      `json:"tags"`
	UserName       string      `json:"user_name"`
}

func (wsh *AnnotatedDataCheckHandler) HandleAnnotatedDataCheck(c echo.Context) error {
	conn, err := wsh.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := wsh.SendNextData(conn, 0); err != nil {
		fmt.Println("[Check]: ", err)
		return nil
	}

	for {
		// Read and Write
		if err := wsh.ReadAndWriteWebSocket(conn); err != nil {
			fmt.Println("[Check]: ", err)
			return nil
		}
	}
}

func SendAnnotatedData(ws *websocket.Conn, fileName string, actPoint, annotatedPoint point.Point, control, userName, tags string) error {
	annotatedData := AnnotateData{
		FileName:       fileName,
		ActPoint:       actPoint,
		AnnotatedPoint: annotatedPoint,
		Control:        control,
		Tags:           tags,
		UserName:       userName,
	}
	data, err := json.Marshal(annotatedData)
	if err != nil {
		return err
	}
	return ws.WriteMessage(websocket.TextMessage, data)
}

func (wsh *AnnotatedDataCheckHandler) ReadAndWriteWebSocket(ws *websocket.Conn) error {
	_, msg, err := ws.ReadMessage()
	if err != nil {
		return err
	}

	var data AnnotateData
	err = json.Unmarshal(msg, &data)
	if err != nil {
		return err
	}
	recivedDataId, _ := carDataDB.SelectIdFromFileName(data.FileName)

	if data.Control == PREV {
		return wsh.SendPrevData(ws, recivedDataId)
	}

	switch data.Control {
	case DELETE:
		carDataDB.DeleteCarData(data.FileName)
	case MOD:
		clickPoint := data.AnnotatedPoint
		magnitude := point.CalcNormalizedMagnitude(wsh.BasePoint, clickPoint, wsh.MaxDistancePoint)
		angle := -(point.CalcAngle(wsh.BasePoint, clickPoint))
		carDataDB.UpdateCarData(data.FileName, magnitude, angle, data.UserName, data.Tags)
	case NEXT:
	}
	return wsh.SendNextData(ws, recivedDataId)
}

func (wsh *AnnotatedDataCheckHandler) SendPrevData(ws *websocket.Conn, id int64) error {
	var carData *carDataDB.CarData
	var err error

	// INFO 最後のページから前ページに戻る際に id == -1 になる
	if id == -1 {
		carData, err = carDataDB.SelectLatestCarData()
	} else {
		carData, err = carDataDB.SelectPrevCarData(id)
	}

	if err == sql.ErrNoRows {
		return SendAnnotatedData(ws, "", point.Point{X: 0, Y: 0}, point.Point{X: 0, Y: 0}, FINISH, "", "")
	}
	if err != nil {
		return err
	}

	actMagnitude := carData.CarSpeed
	actAngle := carData.CarSteering
	actPoint := point.ReverseCalculate(wsh.BasePoint, wsh.MaxDistancePoint, actMagnitude, actAngle)
	annotatedMagnitude := carData.IdealSpeed
	annotatedAngle := carData.IdealSteering
	annotatedPoint := point.ReverseCalculate(wsh.BasePoint, wsh.MaxDistancePoint, annotatedMagnitude, annotatedAngle)

	return SendAnnotatedData(ws, carData.FileName, actPoint, annotatedPoint, NORMAL, carData.AnnotationUserName, carData.Tags)
}

func (wsh *AnnotatedDataCheckHandler) SendNextData(ws *websocket.Conn, id int64) error {
	var carData *carDataDB.CarData
	var err error

	// INFO 最後のページから前ページに戻る際に id == -1 になる
	if id == -1 {
		carData, err = carDataDB.SelectOldestCarData()
	} else {
		carData, err = carDataDB.SelectNextCarData(id)
	}

	if err == sql.ErrNoRows {
		return SendAnnotatedData(ws, "", point.Point{X: 0, Y: 0}, point.Point{X: 0, Y: 0}, FINISH, "", "")
	}
	if err != nil {
		return err
	}

	actMagnitude := carData.CarSpeed
	actAngle := carData.CarSteering
	actPoint := point.ReverseCalculate(wsh.BasePoint, wsh.MaxDistancePoint, actMagnitude, actAngle)
	annotatedMagnitude := carData.IdealSpeed
	annotatedAngle := carData.IdealSteering
	annotatedPoint := point.ReverseCalculate(wsh.BasePoint, wsh.MaxDistancePoint, annotatedMagnitude, annotatedAngle)

	return SendAnnotatedData(ws, carData.FileName, actPoint, annotatedPoint, NORMAL, carData.AnnotationUserName, carData.Tags)
}

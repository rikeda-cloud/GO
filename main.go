package main

import (
	"go_p/db"
	"go_p/handlers"
	"log"

	"github.com/labstack/echo/v4"
)

func initCarData() {
	carDataDB.InsertCarData("1.png", 0.1, -30)
	carDataDB.InsertCarData("2.png", 0.2, -20)
	carDataDB.InsertCarData("3.png", 0.3, -10)
	carDataDB.InsertCarData("4.png", 0.4, 0)
	carDataDB.InsertCarData("5.png", 0.5, 10)
	carDataDB.InsertCarData("6.png", 0.6, 20)
	carDataDB.InsertCarData("7.png", 0.7, 30)
	carDataDB.InsertCarData("8.png", 0.8, 40)
	carDataDB.InsertCarData("9.png", 0.9, 50)
}

func main() {
	if err := carDataDB.CreateCarDataTableIf(); err != nil {
		log.Fatal(err)
	}
	initCarData()
	wsHandler := handlers.NewImageClickHandler()
	e := echo.New()
	e.Static("/", "./static")
	e.Static("/images/", "./images/")
	e.GET("/ws", wsHandler.ImageClickHandler)
	e.Logger.Fatal(e.Start(":8000"))
}

package main

import (
	"GO/cmd/annotation/handlers"
	"GO/internal/config"
	"GO/internal/db"
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
	cfg := config.GetConfig()
	if err := carDataDB.CreateCarDataTableIf(); err != nil {
		log.Fatal(err)
	}
	initCarData()
	imageClickHandler := handlers.NewImageClickHandler()
	remainCountHandler := handlers.NewRemainImageCountHandler()

	e := echo.New()
	e.Static("/", cfg.App.Annotation.StaticDir)
	e.Static("/images/", cfg.Image.DirPath)
	e.GET("/ws", imageClickHandler.ImageClickHandler)
	e.GET("/ws/remain-count", remainCountHandler.RemainImageCountHandler)
	e.Logger.Fatal(e.Start(cfg.App.Annotation.Port))
}

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
	carDataDB.InsertPredictedCarData("5.png", 0.5, 10)
	carDataDB.InsertPredictedCarData("6.png", 0.6, 20)
	carDataDB.InsertPredictedCarData("7.png", 0.7, 30)
	carDataDB.InsertPredictedCarData("8.png", 0.8, 40)
	carDataDB.InsertPredictedCarData("9.png", 0.9, 50)
}

func main() {
	cfg := config.GetConfig()
	if err := carDataDB.CreateCarDataTableIf(); err != nil {
		log.Fatal(err)
	}
	initCarData()
	annotationHandler := handlers.NewAnnotationHandler()
	remainCountHandler := handlers.NewRemainImageCountHandler()
	annotatedDataCheckHandler := handlers.NewAnnotatedDataCheckHandler()
	predictedDataHandler := handlers.NewPredictedDataHandler()
	predictedRemainImageCountHandler := handlers.NewPredictedRemainImageCountHandler()

	e := echo.New()
	e.Static("/", cfg.App.Annotation.StaticDir)
	e.Static("/images/", cfg.Image.DirPath)
	e.Static("/predict-images/", cfg.Image.PredictDirPath)
	e.GET("/ws", annotationHandler.AnnotationHandler)
	e.GET("/ws/remain-count", remainCountHandler.RemainImageCountHandler)
	e.GET("/ws/predict-remain-count", predictedRemainImageCountHandler.PredictedRemainImageCountHandler)
	e.GET("/ws/check", annotatedDataCheckHandler.AnnotatedDataCheckHandler)
	e.GET("/ws/ai", predictedDataHandler.PredictedDataHandler)
	e.Logger.Fatal(e.Start(cfg.App.Annotation.Port))
}

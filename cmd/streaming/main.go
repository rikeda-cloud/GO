package main

import (
	"GO/cmd/streaming/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	imageStreamingHandler := handlers.NewImageStreamingHandler()
	e := echo.New()
	e.Static("/", "./cmd/streaming/static")
	e.GET("/feed", imageStreamingHandler.Handler)
	e.Logger.Fatal(e.Start(":8000"))
}

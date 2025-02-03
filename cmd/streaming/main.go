package main

import (
	"GO/cmd/streaming/handlers"
	"GO/internal/config"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.GetConfig()
	imageStreamingHandler := handlers.NewImageStreamingHandler()
	e := echo.New()
	e.Static("/", cfg.App.Streaming.StaticDir)
	e.GET("/feed", imageStreamingHandler.HandleImageStreaming)
	e.Logger.Fatal(e.Start(cfg.App.Streaming.Port))
}

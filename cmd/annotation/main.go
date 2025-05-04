package main

import (
	"GO/cmd/annotation/router"
	"GO/internal/config"
	"GO/internal/db"
	"log"
)

func main() {
	if err := carDataDB.CreateCarDataTableIf(); err != nil {
		log.Fatal(err)
	}
	carDataDB.InitTmpCarData()
	e := router.SetupRouter()
	cfg := config.GetConfig()
	e.Logger.Fatal(e.Start(cfg.App.Annotation.Port))
}

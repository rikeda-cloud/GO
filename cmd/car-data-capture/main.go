package main

import (
	"GO/cmd/car-data-capture/capture"
	"GO/internal/config"
	"gocv.io/x/gocv"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	cfg := config.GetConfig()
	camera, err := gocv.OpenVideoCapture(cfg.Camera.DeviceNumber)
	if err != nil {
		log.Fatal(err)
	}
	defer camera.Close()

	var speed float64
	var steering float64

	go captureData.CaptureLoop(&wg, camera, &speed, &steering)
	go captureData.DriveTrainingCar(&wg, &speed, &steering)
	wg.Wait()
}

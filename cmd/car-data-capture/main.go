package main

import (
	"GO/cmd/car-data-capture/capture"
	"gocv.io/x/gocv"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	camera, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatal(err)
	}
	defer camera.Close()

	var speed float64
	var steering float64

	go captureData.CaptureLoop(&wg, camera, &speed, &steering)
	wg.Wait()
}

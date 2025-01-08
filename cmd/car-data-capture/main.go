package main

import (
	"GO/cmd/car-data-capture/capture"
)

func main() {
	camera, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatal(err)
	}
	defer camera.Close()

	var speed float64
	var steering float64

	go capture.CaptureLoop(camera, &speed, &steering)
}

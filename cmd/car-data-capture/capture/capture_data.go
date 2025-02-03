package captureData

import (
	"GO/internal/config"
	"GO/internal/db"
	"fmt"
	"log"
	"sync"
	"time"

	"gocv.io/x/gocv"
)

func CaptureLoop(wg *sync.WaitGroup, camera *gocv.VideoCapture, speed, steering *float64) error {
	defer wg.Done()
	img := gocv.NewMat()
	defer img.Close()
	carDataDB.CreateCarDataTableIf()

	cfg := config.GetConfig()
	camera.Set(gocv.VideoCaptureFrameWidth, cfg.Camera.Width)
	camera.Set(gocv.VideoCaptureFrameHeight, cfg.Camera.Height)

	for {
		camera.Read(&img)
		if img.Empty() {
			continue
		}

		fileName := fmt.Sprintf("%d.png", time.Now().UnixNano())
		filePath := cfg.Image.DirPath + fileName
		if ok := gocv.IMWrite(filePath, img); !ok {
			log.Fatal("Error: gocv.IMWrite()")
		}

		sp := *speed
		st := *steering
		carDataDB.InsertCarData(fileName, sp, st)
		fmt.Printf("Capture: %s(%f, %f)\n", fileName, sp, st)

		time.Sleep(cfg.App.CarDataCapture.CaptureIntervalMsec * time.Millisecond)
	}
}

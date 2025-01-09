package captureData

import (
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

	camera.Set(gocv.VideoCaptureFrameWidth, 640)
	camera.Set(gocv.VideoCaptureFrameHeight, 480)

	for {
		camera.Read(&img)
		if img.Empty() {
			continue
		}

		fileName := fmt.Sprintf("%d.png", time.Now().UnixNano())
		filePath := "./images/" + fileName
		if ok := gocv.IMWrite(filePath, img); !ok {
			log.Fatal("Error: gocv.IMWrite()")
		}

		carDataDB.InsertCarData(fileName, *speed, *steering)

		time.Sleep(100 * time.Millisecond)
	}
}

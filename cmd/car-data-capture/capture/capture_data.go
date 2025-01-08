package captureData

import (
	"GO/internal/db"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	"gocv.io/x/gocv"
)

func CaptureLoop(wg *sync.WaitGroup, camera *gocv.VideoCapture, speed, steering *float64) error {
	defer wg.Done()
	img := gocv.NewMat()
	defer img.Close()
	carDataDB.CreateCarDataTableIf()

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

		// TODO Change Actual speed & steering
		*speed = rand.Float64()
		*steering = float64(rand.IntN(181) - 90)
		carDataDB.InsertCarData(fileName, *speed, *steering)

		time.Sleep(1 * time.Second)
	}
}

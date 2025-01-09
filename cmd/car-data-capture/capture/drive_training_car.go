package captureData

import (
	"math/rand"
	"sync"
	"time"
)

func DriveTrainingCar(wg *sync.WaitGroup, speed, steering *float64) error {
	defer wg.Done()

	for {
		*speed = rand.Float64()
		*steering = float64(rand.Intn(181) - 90)
		time.Sleep(time.Millisecond * 50)
	}
}

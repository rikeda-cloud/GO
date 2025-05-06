package exporter

import (
	"GO/internal/db"
	"fmt"
)

func ExportToStdout(data *carDataDB.CarData, values []float64) error {
	fmt.Println(data.IdealSpeed, data.IdealSteering, data.CreatedAt, values)
	return nil
}

package exporter

import (
	"encoding/json"
	"io"
	"os"

	"GO/internal/config"
	"GO/internal/db"
	"fmt"
)

type ExportData struct {
	IdealSpeed    float64
	IdealSteering float64
	CreatedAt     string
	Features      []float64
}

func ExportToJson(data *carDataDB.CarData, values []float64) error {
	cfg := config.GetConfig()

	file, err := os.OpenFile(cfg.App.DataExporter.JsonFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("faild to open json file: %w", err)
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)

	var exportDataSlice []ExportData
	if len(byteValue) != 0 {
		if err := json.Unmarshal(byteValue, &exportDataSlice); err != nil {
			return fmt.Errorf("faild to unmarshal json data: %w", err)
		}
	}

	newExportData := ExportData{
		IdealSpeed:    data.IdealSpeed,
		IdealSteering: data.IdealSteering,
		CreatedAt:     data.CreatedAt,
		Features:      values,
	}

	exportDataSlice = append(exportDataSlice, newExportData)
	file.Truncate(0)
	file.Seek(0, 0)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportDataSlice); err != nil {
		return fmt.Errorf("faild to encode: %w", err)
	}
	return nil
}

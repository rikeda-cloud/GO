package exporter

import (
	"GO/internal/config"
	"GO/internal/db"
)

type ExportFunc func(*carDataDB.CarData, []float64) error

func GetExportFunc() ExportFunc {
	cfg := config.GetConfig()

	switch cfg.App.DataExporter.ExportTo {
	case "cloud":
		return ExportToCloud
	case "json":
		return ExportToJson
	default:
		return ExportToStdout
	}
}

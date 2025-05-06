package exporter

import (
	"bytes"
	"encoding/json"
	"net/http"

	"GO/internal/config"
	"GO/internal/db"
	"fmt"
)

func ExportToCloud(data *carDataDB.CarData, values []float64) error {
	cfg := config.GetConfig()

	content := fmt.Sprintf(
		"**走行データの送信**\n"+
			"IdealSpeed: %.2f\n"+
			"IdealSteering: %.2f\n"+
			"CreatedAt: %s\n"+
			"Features: %v",
		data.IdealSpeed,
		data.IdealSteering,
		data.CreatedAt,
		values,
	)

	message := map[string]string{
		"content": content,
	}

	jsonBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(cfg.App.DataExporter.CloudURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to POST: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("cloud response error: %s", resp.Status)
	}

	return nil
}

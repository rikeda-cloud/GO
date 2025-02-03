package handlers

import (
	"GO/internal/config"
	"GO/internal/point"
)

type CoordinateRange struct {
	BasePoint        point.Point
	MaxDistancePoint point.Point
}

func NewCoordinateRange() *CoordinateRange {
	cfg := config.GetConfig()
	return &CoordinateRange{
		BasePoint:        point.Point{X: float64(cfg.Camera.Width / 2), Y: float64(cfg.Camera.Height)},
		MaxDistancePoint: point.Point{X: 0, Y: cfg.Camera.Height},
	}
}

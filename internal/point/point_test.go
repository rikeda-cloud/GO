package point

import (
	"testing"
)

func TestCalcAngle(t *testing.T) {
	tests := []struct {
		x      float64
		y      float64
		expect float64
	}{
		{x: 0, y: 480, expect: -90},
		{x: 640, y: 480, expect: 90},
		{x: 160, y: 320, expect: -45},
		{x: 320, y: 0, expect: 0},
	}
	basePoint := Point{X: 320, Y: 480}

	for _, tc := range tests {
		clickPoint := Point{X: tc.x, Y: tc.y}
		angle := -(CalcAngle(basePoint, clickPoint))
		if angle != tc.expect {
			t.Fatalf("expect: %f, act: %f", tc.expect, angle)
		}
	}
}

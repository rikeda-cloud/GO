package point

import (
	"math"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func CalcDistance(p1, p2 Point) float64 {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func CalcAngle(p1, p2 Point) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return math.Atan2(dx, dy) * 180 / math.Pi
}

func CalcNormalizedMagnitude(basePoint, clickPoint, maxDistancePoint Point) float64 {
	distance := CalcDistance(basePoint, clickPoint)
	maxDistance := CalcDistance(basePoint, maxDistancePoint)
	normalizedMagnitude := distance / maxDistance
	return normalizedMagnitude
}

func ReverseCalculate(basePoint, maxDistancePoint Point, magnitude, angle float64) Point {
	maxDistance := CalcDistance(basePoint, maxDistancePoint)
	distance := magnitude * maxDistance

	angleRad := angle * math.Pi / 180

	// 方向ベクトル
	deltaX := distance * math.Sin(angleRad)
	deltaY := distance * math.Cos(angleRad)

	// 座標
	originalX := basePoint.X + deltaX
	originalY := basePoint.Y - deltaY

	return Point{X: originalX, Y: originalY}
}

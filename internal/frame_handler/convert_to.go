package frameHandler

import (
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

func ConvertToGray(frame *gocv.Mat) *gocv.Mat {
	grayFrame := gocv.NewMat()
	gocv.CvtColor(frame, &grayFrame, gocv.ColorBGRToGray)
	return grayFrame
}

func ConvertToCanny(frame *gocv.Mat) *gocv.Mat {
	grayFrame := ConvertToGray(frame)
	cannyFrame := gocv.NewMat()
	gocv.Canny(
		grayFrame,
		&cannyFrame,
		100,
		200,
	)
	return cannyFrame
}

func ConvertToHough(frame *gocv.Mat) *gocv.Mat {
	houghFrame := gocv.NewMat()
	frame.CopyTo(&houghFrame)
	cannyFrame := ConvertToCanny(frame)
	lines := gocv.NewMat()
	gocv.HoughLinesPWithParams(
		cannyFrame,
		&lines,
		8.0,
		math.Pi/60.0,
		100,
		100.0,
		5.0,
	)
	red := color.RGBA{R: 255, G: 0, B: 0, A: 0}
	for i := 0; i < lines.Rows(); i++ {
		line := lines.GetVeciAt(i, 0)
		pt1 := image.Pt(int(line[0]), int(line[1]))
		pt2 := image.Pt(int(line[2]), int(line[3]))
		gocv.Line(&houghFrame, pt1, pt2, red, 2)
	}
	return houghFrame
}

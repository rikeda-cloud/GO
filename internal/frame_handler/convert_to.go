package frameHandler

import (
	"GO/internal/config"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

func ConvertToGray(frame *gocv.Mat) gocv.Mat {
	grayFrame := gocv.NewMat()
	gocv.CvtColor(*frame, &grayFrame, gocv.ColorBGRToGray)
	return grayFrame
}

func ConvertToCanny(frame *gocv.Mat) gocv.Mat {
	cfg := config.GetConfig()
	grayFrame := ConvertToGray(frame)
	defer grayFrame.Close()
	cannyFrame := gocv.NewMat()
	gocv.Canny(
		grayFrame,
		&cannyFrame,
		cfg.Frame.Canny.Threshold1,
		cfg.Frame.Canny.Threshold2,
	)
	return cannyFrame
}

func DetectHoughData(frame *gocv.Mat) gocv.Mat {
	cfg := config.GetConfig()
	cannyFrame := ConvertToCanny(frame)
	defer cannyFrame.Close()

	lines := gocv.NewMat()

	gocv.HoughLinesPWithParams(
		cannyFrame,
		&lines,
		cfg.Frame.Hough.Rho,
		math.Pi/cfg.Frame.Hough.Step,
		cfg.Frame.Hough.Threshold,
		cfg.Frame.Hough.MinLineLength,
		cfg.Frame.Hough.MaxLineGap,
	)
	return lines
}

func ConvertToHough(frame *gocv.Mat) gocv.Mat {
	houghFrame := gocv.NewMat()
	frame.CopyTo(&houghFrame)

	lines := DetectHoughData(frame)
	defer lines.Close()

	red := color.RGBA{R: 255, G: 0, B: 0, A: 0}
	for i := 0; i < lines.Rows(); i++ {
		line := lines.GetVeciAt(i, 0)
		pt1 := image.Pt(int(line[0]), int(line[1]))
		pt2 := image.Pt(int(line[2]), int(line[3]))
		gocv.Line(&houghFrame, pt1, pt2, red, 2)
	}
	return houghFrame
}

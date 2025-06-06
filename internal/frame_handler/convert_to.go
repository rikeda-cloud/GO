package frameHandler

import (
	"GO/internal/config"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

// INFO Grayスケール画像に変換
func ConvertToGray(frame *gocv.Mat) gocv.Mat {
	grayFrame := gocv.NewMat()
	gocv.CvtColor(*frame, &grayFrame, gocv.ColorBGRToGray)
	return grayFrame
}

// INFO Cannyエッジ検出後の画像に変換
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

// INFO Hough変換で直線データを取得
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

// INFO Hough変換を用いて直線を強調した画像に変換
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

// INFO 左右反転画像に変換
func ConvertToReverse(frame *gocv.Mat) gocv.Mat {
	reversedFrame := gocv.NewMat()
	gocv.Flip(*frame, &reversedFrame, 1)
	return reversedFrame
}

// INFO ぼかし(ノイズ除去。エッジ検出と併用可能)
func ConvertToBilateralFilter(frame *gocv.Mat) gocv.Mat {
	cfg := config.GetConfig()
	if IsGrayscale(frame) {
		return frame.Clone()
	}

	filteredFrame := gocv.NewMat()
	gocv.BilateralFilter(*frame, &filteredFrame, cfg.Frame.Filter.D, cfg.Frame.Filter.SigmaColor, cfg.Frame.Filter.SigmaSpace)
	return filteredFrame
}

// INFO 画像の白黒化
func ConvertToBinary(frame *gocv.Mat) gocv.Mat {
	cfg := config.GetConfig()

	binary := gocv.NewMat()
	grayFrame := ConvertToGray(frame)
	defer grayFrame.Close()

	gocv.Threshold(grayFrame, &binary, cfg.Frame.Binary.Threshold, cfg.Frame.Binary.MaxValue, gocv.ThresholdBinary)

	return binary
}

// INFO HaarLike特徴量抽出を行い、画像にプロット
func ConvertToHaarLike(frame *gocv.Mat) gocv.Mat {
	cfg := config.GetConfig()
	haarValues := CalcHaarValues(frame, cfg.Frame.HaarLike.Divisions, cfg.Frame.HaarLike.RectHeight)

	width := frame.Cols()
	height := frame.Rows()
	widthStep := width / cfg.Frame.HaarLike.Divisions

	resultFrame := frame.Clone()
	black := color.RGBA{0, 0, 0, 0}
	for i, val := range haarValues {
		x := i * widthStep
		y := int(val * float64(height))
		rect := image.Rect(x, y, x+widthStep, y+1)
		gocv.Rectangle(&resultFrame, rect, black, 1)
	}

	return resultFrame
}

// INFO SIFT特徴量抽出を行い、画像にプロット
func ConvertToSift(frame *gocv.Mat) gocv.Mat {
	sift := gocv.NewSIFT()
	defer sift.Close()

	kps, _ := sift.DetectAndCompute(*frame, gocv.NewMat())

	// 特徴点が1つもなければそのまま返す
	if len(kps) == 0 {
		return *frame
	}

	siftFrame := gocv.NewMat()
	gocv.DrawKeyPoints(*frame, kps, &siftFrame, color.RGBA{R: 0, G: 255, B: 0, A: 0}, 4)

	return siftFrame
}

// INFO AKAZE特徴量抽出を行い、画像にプロット
func ConvertToAKAZE(frame *gocv.Mat) gocv.Mat {
	akaze := gocv.NewAKAZE()
	defer akaze.Close()

	keypoints, _ := akaze.DetectAndCompute(*frame, gocv.NewMat())

	akazeFrame := frame.Clone()
	for _, kp := range keypoints {
		pt := image.Pt(int(kp.X), int(kp.Y))
		gocv.Circle(&akazeFrame, pt, 2, color.RGBA{255, 0, 0, 0}, 2)
	}
	return akazeFrame
}

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

func ConvertToHaarLike(frame *gocv.Mat) gocv.Mat {
	DIVISIONS := 40
	RECT_HEIGHT := 15
	// グレースケール変換
	grayFrame := ConvertToGray(frame)

	width := grayFrame.Cols()
	height := grayFrame.Rows()
	widthStep := width / DIVISIONS

	// 各領域のHaar-like特徴量を計算
	haarValues := make([]float64, DIVISIONS)
	for i := 0; i < DIVISIONS; i++ {
		x := i * widthStep
		roi := grayFrame.Region(image.Rect(x, 0, x+widthStep, height))
		defer roi.Close()

		// 各行の輝度の合計を計算
		rowSums := make([]float64, roi.Rows())
		for y := 0; y < roi.Rows(); y++ {
			sum := 0.0
			for x := 0; x < roi.Cols(); x++ {
				sum += float64(roi.GetUCharAt(y, x))
			}
			rowSums[y] = sum
		}

		// 平均化カーネルを作成
		kernel := make([]float64, RECT_HEIGHT)
		for k := range kernel {
			kernel[k] = 1.0 / float64(RECT_HEIGHT)
		}

		// 畳み込みを実行
		convolved := convolve(rowSums, kernel)

		// 差分を計算
		diff := make([]float64, len(convolved)-1)
		for j := 0; j < len(convolved)-1; j++ {
			diff[j] = convolved[j+1] - convolved[j]
		}

		// 最大の差分のインデックスを取得
		maxIdx := 0
		maxVal := diff[0]
		for j, val := range diff {
			if val > maxVal {
				maxVal = val
				maxIdx = j
			}
		}

		// 正規化して保存
		haarValues[i] = float64(maxIdx) / float64(len(diff))
	}

	// 元の画像に矩形を描画
	result := frame.Clone()
	black := color.RGBA{0, 0, 0, 0}
	for i, val := range haarValues {
		x := i * widthStep
		y := int(val * float64(height))
		rect := image.Rect(x, y, x+widthStep, y+1)
		gocv.Rectangle(&result, rect, black, 1)
	}

	return result
}

// 1次元配列に対する畳み込み関数
func convolve(array []float64, kernel []float64) []float64 {
	kernelSize := len(kernel)
	arrayLen := len(array)
	resultLen := arrayLen - kernelSize + 1
	result := make([]float64, resultLen)

	for i := 0; i < resultLen; i++ {
		sum := 0.0
		for j := 0; j < kernelSize; j++ {
			sum += array[i+j] * kernel[kernelSize-1-j]
		}
		result[i] = sum
	}

	return result
}

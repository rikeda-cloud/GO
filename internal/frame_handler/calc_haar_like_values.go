package frameHandler

import (
	"image"
	"sync"

	"gocv.io/x/gocv"
)

// HaarLike特徴量抽出で画像の白黒差が激しいポイントを算出する
func CalcHaarValues(frame *gocv.Mat, divisions, rectHeight int) []float64 {
	grayFrame := ConvertToGray(frame)

	width := grayFrame.Cols()
	height := grayFrame.Rows()
	widthStep := width / divisions

	haarValues := make([]float64, divisions)
	var wg sync.WaitGroup

	// 並列に各divisionを処理
	for i := 0; i < divisions; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			x := i * widthStep
			roi := grayFrame.Region(image.Rect(x, 0, x+widthStep, height))
			defer roi.Close()

			rowSums := make([]float64, roi.Rows())
			for y := 0; y < roi.Rows(); y++ {
				sum := 0.0
				for x := 0; x < roi.Cols(); x++ {
					sum += float64(roi.GetUCharAt(y, x))
				}
				rowSums[y] = sum
			}

			kernel := make([]float64, rectHeight)
			for k := range kernel {
				kernel[k] = 1.0 / float64(rectHeight)
			}

			convolved := convolve(rowSums, kernel)

			diff := make([]float64, len(convolved)-1)
			for j := 0; j < len(convolved)-1; j++ {
				diff[j] = convolved[j+1] - convolved[j]
			}

			maxIdx := 0
			maxVal := diff[0]
			for j, val := range diff {
				if val > maxVal {
					maxVal = val
					maxIdx = j
				}
			}

			haarValues[i] = float64(maxIdx) / float64(len(diff))
		}(i)
	}

	wg.Wait()
	return haarValues
}

// 1次元配列に対する畳み込み関数
func convolve(array, kernel []float64) []float64 {
	kernelSize := len(kernel)
	arrayLen := len(array)
	resultLen := arrayLen - kernelSize + 1
	result := make([]float64, resultLen)

	// あらかじめカーネルを反転
	reversedKernel := make([]float64, kernelSize)
	for i := 0; i < kernelSize; i++ {
		reversedKernel[i] = kernel[kernelSize-1-i]
	}

	for i := 0; i < resultLen; i++ {
		sum := 0.0
		for j := 0; j < kernelSize; j++ {
			sum += array[i+j] * reversedKernel[j]
		}
		result[i] = sum
	}
	return result
}

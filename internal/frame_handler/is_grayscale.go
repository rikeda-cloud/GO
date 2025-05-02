package frameHandler

import (
	"gocv.io/x/gocv"
)

func IsGrayscale(mat *gocv.Mat) bool {
	return mat.Channels() == 1
}

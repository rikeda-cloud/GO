package handlers

import (
	"log"
	"time"

	"GO/internal/frame_handler"
	"github.com/labstack/echo/v4"
	"gocv.io/x/gocv"
)

type ImageStreamingHandler struct {
	camera *gocv.VideoCapture
}

func NewImageStreamingHandler() *ImageStreamingHandler {
	camera, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatal(err)
	}
	camera.Set(gocv.VideoCaptureFrameWidth, 640)
	camera.Set(gocv.VideoCaptureFrameHeight, 480)

	return &ImageStreamingHandler{
		camera,
	}
}

func (wsh *ImageStreamingHandler) Handler(c echo.Context) error {
	defer wsh.camera.Close()

	img := gocv.NewMat()
	defer img.Close()

	c.Response().Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")

	for {
		// フレームレート調整 (30FPSの場合: 33msのスリープ)
		// time.Sleep(33 * time.Millisecond)

		if ok := wsh.camera.Read(&img); !ok || img.Empty() {
			log.Println("Error Capture Image")
			continue
		}

		start := time.Now()
		houghImg := frameHandler.ConvertToHough(&img)
		log.Println(time.Since(start))

		buf, err := gocv.IMEncode(".png", houghImg)
		if err != nil {
			log.Println(err)
		}

		// フレームをストリームとして送信
		if _, err := c.Response().Write([]byte("--frame\r\n")); err != nil {
			continue
		}
		if _, err := c.Response().Write([]byte("Content-Type: image/jpeg\r\n\r\n")); err != nil {
			continue
		}
		if _, err := c.Response().Write(buf.GetBytes()); err != nil {
			continue
		}
		if _, err := c.Response().Write([]byte("\r\n")); err != nil {
			continue
		}
		c.Response().Flush()
	}
}

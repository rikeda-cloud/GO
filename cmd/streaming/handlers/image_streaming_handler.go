package handlers

import (
	"log"
	"time"

	"GO/internal/config"
	"GO/internal/frame_handler"
	"GO/internal/ws"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gocv.io/x/gocv"
)

type ImageStreamingHandler struct {
	ws.WebSocketBaseHandler
	camera          *gocv.VideoCapture
	image_converter func(*gocv.Mat) gocv.Mat
}

func NewImageStreamingHandler() *ImageStreamingHandler {
	cfg := config.GetConfig()
	camera, err := gocv.OpenVideoCapture(cfg.Camera.DeviceNumber)
	if err != nil {
		log.Fatal(err)
	}
	camera.Set(gocv.VideoCaptureFrameWidth, cfg.Camera.Width)
	camera.Set(gocv.VideoCaptureFrameHeight, cfg.Camera.Height)

	return &ImageStreamingHandler{
		*ws.NewWebSocketBaseHandler(),
		camera,
		frameHandler.ConvertToHough, // デフォルトではハフ変換を適用
	}
}

func (wsh *ImageStreamingHandler) HandleImageStreaming(c echo.Context) error {
	conn, err := wsh.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer wsh.camera.Close()

	img := gocv.NewMat()
	defer img.Close()

	for {
		// カメラからフレームを取得
		if ok := wsh.camera.Read(&img); !ok || img.Empty() {
			log.Println("Error capturing image")
			continue
		}

		start := time.Now()
		convertedImg := wsh.image_converter(&img)
		log.Println(time.Since(start))

		buf, err := gocv.IMEncode(".png", convertedImg)
		if err != nil {
			log.Println(err)
			continue
		}

		// WebSocketでエンコードされた画像を送信
		if err := conn.WriteMessage(websocket.BinaryMessage, buf.GetBytes()); err != nil {
			log.Println("WebSocket Write Error:", err)
			break
		}
		time.Sleep(30 * time.Millisecond)
	}
	return nil
}

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
	cfg := config.GetConfig()

	img := gocv.NewMat()
	defer img.Close()

	done := make(chan struct{})

	// クライアントからのメッセージ受信をgoroutineで処理
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				close(done)
				break
			}
			switch string(msg) {
			case "1":
				wsh.image_converter = frameHandler.ConvertToHough
			case "2":
				wsh.image_converter = frameHandler.ConvertToGray
			case "3":
				wsh.image_converter = frameHandler.ConvertToCanny
			default:
				wsh.image_converter = frameHandler.ConvertToHough
			}
		}
	}()

	for {
		// ゴルーチンが終了したかを待つ
		select {
		case <-done:
			// ゴルーチンが終了
			return nil
		default:
			// ゴルーチンが終了していない場合は何もしない
		}

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
		time.Sleep(cfg.App.Streaming.StreamingIntervalMsec * time.Millisecond)
	}
	return nil
}

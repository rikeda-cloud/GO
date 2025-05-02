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
	camera        *gocv.VideoCapture
	frame_handler func(*gocv.Mat) gocv.Mat
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
	// defer wsh.camera.Close()

	cfg := config.GetConfig()
	img := gocv.NewMat()
	defer img.Close()

	done := make(chan struct{})

	// クライアントからのメッセージを用いて画像処理関数を切り替える
	go wsh.ReadFromWebSocket(conn, done)

	for {
		select {
		case <-done: // ゴルーチンが終了
			return nil
		default:
		}

		// カメラからフレームを取得
		if ok := wsh.camera.Read(&img); !ok || img.Empty() {
			log.Println("Error capturing image")
			continue
		}

		start := time.Now()
		convertedImg := wsh.frame_handler(&img)
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

func (wsh *ImageStreamingHandler) ReadFromWebSocket(ws *websocket.Conn, done chan struct{}) {
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			close(done)
			return
		}
		wsh.ChangeFrameHandler(string(msg))
	}
}

func (wsh *ImageStreamingHandler) ChangeFrameHandler(handler_number string) {
	switch handler_number {
	case "1":
		wsh.frame_handler = frameHandler.ConvertToHough
	case "2":
		wsh.frame_handler = frameHandler.ConvertToGray
	case "3":
		wsh.frame_handler = frameHandler.ConvertToCanny
	case "4":
		wsh.frame_handler = frameHandler.ConvertToReverse
	case "5":
		wsh.frame_handler = frameHandler.ConvertToBilateralFilter
	case "6":
		wsh.frame_handler = frameHandler.ConvertToBinary
	case "7":
		wsh.frame_handler = frameHandler.ConvertToHaarLike
	default:
		wsh.frame_handler = frameHandler.ConvertToHough
	}
}

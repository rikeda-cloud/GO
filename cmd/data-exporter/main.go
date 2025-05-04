package main

import (
	"GO/cmd/data-exporter/exporter"
	"GO/internal/config"
	"GO/internal/db"
	"GO/internal/frame_handler"
	"fmt"
	"time"

	"gocv.io/x/gocv"
)

func main() {
	cfg := config.GetConfig()
	id := 0 // INFO 最小のidを初期値として使用

	total := 0 // INFO 送信したデータレコード数
	for {
		data, err := carDataDB.SelectNextCarData(int64(id))
		if err != nil {
			fmt.Println("Finish Export Data:", total, "(data)")
			break
		}
		filePath := cfg.Image.DirPath + data.FileName
		img := gocv.IMRead(filePath, gocv.IMReadColor)
		if img.Empty() {
			fmt.Printf("画像の読み込みに失敗: %s\n", filePath)
			img.Close()
			break
		}

		values := frameHandler.CalcHaarValues(&img, cfg.Frame.HaarLike.Divisions, cfg.Frame.HaarLike.RectHeight)
		img.Close()

		if err := exporter.ExportToCloud(data, values); err != nil {
			fmt.Printf("データのExportに失敗: %s\n", err)
			break
		}

		// INFO 一定時間遅延
		time.Sleep(cfg.App.DataExporter.ExportDelayMsec * time.Millisecond)

		id = data.ID
		total++
	}
}

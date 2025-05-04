package main

import (
	"GO/internal/config"
	"GO/internal/db"
	"GO/internal/frame_handler"
	"fmt"
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

		values := frameHandler.CalcHaarValues(&img, 20, 15)
		img.Close()

		fmt.Println(data.IdealSpeed, data.IdealSteering, data.CreatedAt, values)

		id = data.ID
		total++
	}
}

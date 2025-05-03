package config

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Config struct {
	App struct {
		Annotation struct {
			StaticDir string `json:"static_dir"`
			Port      string `json:"port"`
		} `json:"annotation"`
		Streaming struct {
			StaticDir             string        `json:"static_dir"`
			Port                  string        `json:"port"`
			StreamingIntervalMsec time.Duration `json:"streaming_interval_msec"`
		} `json:"streaming"`
		CarDataCapture struct {
			CaptureIntervalMsec time.Duration `json:"capture_interval_msec"`
		} `json:"car-data-capture"`
		IpNotify struct {
			DiscordWebhookUrl string `json:"discord_webhook_url"`
		} `json:"ip-notify"`
	} `json:"app"`
	Database struct {
		DBMS     string `json:"dbms"`
		FilePath string `json:"file_path"`
	} `json:"database"`
	Image struct {
		DirPath        string `json:"dir_path"`
		PredictDirPath string `json:"predict_dir_path"`
	} `json:"image"`
	Camera struct {
		DeviceNumber int     `json:"device_number"`
		Width        float64 `json:"width"`
		Height       float64 `json:"height"`
	} `json:"camera"`
	Frame struct {
		Canny struct {
			Threshold1 float32 `json:"threshold1"`
			Threshold2 float32 `json:"threshold2"`
		} `json:"canny"`
		Hough struct {
			Rho           float32 `json:"rho"`
			Step          float32 `json:"step"`
			Threshold     int     `json:"threshold"`
			MinLineLength float32 `json:"min_line_length"`
			MaxLineGap    float32 `json:"max_line_gap"`
		} `json:"hough"`
		Filter struct {
			D          int     `json:"d"`
			SigmaColor float64 `json:"sigma_color"`
			SigmaSpace float64 `json:"sigma_space"`
		} `json:"filter"`
		Binary struct {
			Threshold float32 `json:"threshold"`
			MaxValue  float32 `json:"max_value"`
		} `json:"binary"`
		HaarLike struct {
			Divisions  int `json:"divisions"`
			RectHeight int `json:"rect_height"`
		} `json:"haar_like"`
	} `json:"frame"`
}

var (
	instance *Config
	once     sync.Once
)

func init() {
	loadConfig("configs/config.json")
}

func loadConfig(filePath string) {
	once.Do(func() {
		data, err := os.ReadFile(filePath)
		if err != nil {
			panic("Failed to read config file: " + err.Error())
		}

		var cfg Config
		if err := json.Unmarshal(data, &cfg); err != nil {
			panic("Failed to parse config file: " + err.Error())
		}

		instance = &cfg
	})
}

func GetConfig() *Config {
	if instance == nil {
		panic("Config is not initialized. Call LoadConfig() first.")
	}
	return instance
}

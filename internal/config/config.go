package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	App struct {
		Annotation struct {
			StaticDir string `json:"static_dir"`
			Port      string `json:"port"`
		} `json:"annotation"`
		Streaming struct {
			StaticDir string `json:"static_dir"`
			Port      string `json:"port"`
		} `json:"streaming"`
	} `json:"app"`
	Database struct {
		DBMS     string `json:"dbms"`
		FilePath string `json:"file_path"`
	} `json:"database"`
	Image struct {
		DirPath string `json:"dir_path"`
	} `json:"image"`
	Camera struct {
		DeviceNumber int `json:"device_number"`
		Width        int `json:"width"`
		Height       int `json:"height"`
	} `json:"camera"`
	Frame struct {
		Canny struct {
			Threshold1 int `json:"threshold1"`
			Threshold2 int `json:"threshold2"`
		} `json:"canny"`
		Hough struct {
			Rho           float64 `json:"rho"`
			Step          float64 `json:"step"`
			Threshold     int64   `json:"threshold"`
			MinLineLength float64 `json:"min_line_length"`
			MaxLineGap    float64 `json:"max_line_gap"`
		} `json:"hough"`
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

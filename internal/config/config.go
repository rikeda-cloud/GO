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

package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	StorageData string      `json:"storage_data"`
	TelegramBot TelegramBot `json:"telegram_bot"`
	Watcher     Watcher     `json:"watcher"`
}
type TelegramBot struct {
	URL string `json:"url"`
}

type Watcher struct {
	IntervalMinutes  int `json:"interval_minutes"`
	MaxReposPerCheck int `json:"max_repos_per_check"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		log.Fatal("Config file path is not exist")
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Config file does not exist: %s", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Cannot read config file: %s", err)
	}

	return &config
}

package config

import (
	"encoding/json"
	"log"
	"os"
)

type KafkaConfig struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type Config struct {
	Token       string       `json:"token"`
	Commands    []BotCommand `json:"commands"`
	ServerURL   string       `json:"server_url"`
	KafkaConfig KafkaConfig  `json:"kafka"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		log.Fatal("Config file path is not exist")
	}
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Config file does not exist:%s", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Cannot read config file:%s", err)
	}

	return &config
}

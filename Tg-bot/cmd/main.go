package main

import (
	"tg-bot/internal/app"
	"tg-bot/internal/config"
)

const filePath = "configs/config.json" // or switch to $PATH

func main() {
	cfg := config.MustLoad(filePath)
	application := app.New(cfg)
	application.Start()
}

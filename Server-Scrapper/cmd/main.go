package main

import (
	"server-scrapper/internal/app"
	"server-scrapper/internal/config"
)

// TODO: ДОБАВИТЬ ТАЙМАУТЫ И ГОРУТИНЫ
const ConfigPath = "configs/bd_config.json"

func main() {
	cfg := config.MustLoad(ConfigPath)
	app := app.New(cfg)
	app.Start()
}

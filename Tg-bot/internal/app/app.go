package app

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"tg-bot/internal/config"
	"tg-bot/internal/lib/network"
	"tg-bot/internal/server"
	"tg-bot/internal/tgbot"
)

type App struct {
	tgbot   *tgbot.TgBot
	network *network.Network
	server  *server.Server
}

func New(cfg *config.Config) *App {
	net := network.New(cfg.ServerURL)
	bot := tgbot.New(cfg)
	server := server.New(bot)
	return &App{tgbot: bot, network: net, server: server}
}

func (app *App) Start() {

	go func() {
		if err := app.server.Start(app.server.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to start server", "error", err)
		}
	}()

	app.tgbot.StopChan = make(chan struct{})
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := app.tgbot.GetUpdatesChan(u)

	go tgbot.MainController(updates, app.tgbot, app.tgbot.StopChan, app.network)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	app.tgbot.StopReceivingUpdates()
	close(app.tgbot.StopChan)
	log.Println("Shutting down tgbot server...")
}

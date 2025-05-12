package app

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/signal"
	"tg-bot/internal/config"
	"tg-bot/internal/lib/kafka"
	"tg-bot/internal/lib/network"
	"tg-bot/internal/server"
	"tg-bot/internal/tgbot"
)

type App struct {
	tgbot   *tgbot.TgBot
	network *network.Network
	server  *server.Server
	kafka   *kafka.Kafka
}

func New(cfg *config.Config) *App {
	net := network.New(cfg.ServerURL)
	bot := tgbot.New(cfg)
	server := server.New(bot)
	kafka := kafka.New(cfg.KafkaConfig, bot)
	return &App{tgbot: bot, network: net, server: server, kafka: kafka}
}

func (app *App) Start() {

	app.tgbot.StopChan = make(chan struct{})
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := app.tgbot.GetUpdatesChan(u)

	go tgbot.MainController(updates, app.tgbot, app.tgbot.StopChan, app.network)

	go app.kafka.ReadMessages()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	app.tgbot.StopReceivingUpdates()
	close(app.tgbot.StopChan)
	log.Println("Shutting down tgbot server...")
}

func (app *App) consumeKafkaMessages() {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	for {
		msg, err := app.kafka.Consumer.ReadMessage(context.Background())
		if err != nil {
			//log.Println("Kafka consumer error:", err)
			continue
		}

		log.Printf("Kafka received: %s\n", string(msg.Value))

		// здесь можно переслать в Telegram или вызвать handler
	}
}

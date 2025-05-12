package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"tg-bot/internal/config"
	"tg-bot/internal/lib/api"
	"tg-bot/internal/tgbot"
)

type Kafka struct {
	Consumer *kafka.Reader
	tgBot    *tgbot.TgBot
	topic    string
}

func New(cfg config.KafkaConfig, bot *tgbot.TgBot) *Kafka {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
	})

	return &Kafka{
		Consumer: consumer,
		tgBot:    bot,
		topic:    cfg.Topic,
	}
}

func (k *Kafka) ReadMessages() {
	for {
		msg, err := k.Consumer.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error while reading message: %v", err)
			continue
		}

		var notification api.GitRepo
		err = json.Unmarshal(msg.Value, &notification)
		if err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		k.tgBot.SendMessage(notification.To, notification.Msg)

	}

}

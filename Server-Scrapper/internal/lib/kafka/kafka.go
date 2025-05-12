package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"server-scrapper/internal/config"
	"server-scrapper/internal/lib/api"
	"server-scrapper/internal/lib/repositories"
	"time"
)

type KafkaWorker interface {
	SendRepoUpdates(repo repositories.WatchedRepo, message string)
}

type Kafka struct {
	producer *kafka.Writer
	topic    string
}

func New(cfg config.Kafka) *Kafka {
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
	})

	return &Kafka{
		producer: producer,
		topic:    cfg.Topic,
	}
}

func (k *Kafka) SendRepoUpdates(repo repositories.WatchedRepo, message string) {
	data := api.GitRepo{To: repo.ChatID, Name: repo.RepoName, Owner: repo.RepoOwner, Msg: message}
	msg, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal message: %v", err)

	}

	kafkaMessage := kafka.Message{
		Value: msg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = k.producer.WriteMessages(ctx, kafkaMessage)
	if err != nil {
		log.Printf("failed to send Kafka message: %v", err)
	}
}

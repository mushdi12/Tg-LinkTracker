package tgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
	"tg-bot/internal/config"
	"tg-bot/internal/lib/convertor"
)

type Bot interface {
	SendMessage(chatId int64, message string)
	GetCommand() []tgbotapi.BotCommand
}
type TgBot struct {
	*tgbotapi.BotAPI
	ServerURL string
	StopChan  chan struct{}
	wg        sync.WaitGroup
}

func New(cfg *config.Config) *TgBot {

	tgBot := &TgBot{}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal("Telegram token is empty")
	}

	_, err = bot.Request(tgbotapi.DeleteWebhookConfig{})
	if err != nil {
		log.Fatalf("Failed to delete webhook: %v", err)
	}

	botCommands, err := convertor.CommandConverter(cfg.Commands)
	if err != nil {
		log.Fatal("Failed to convert commands: ", err)
	}

	setCommands := tgbotapi.NewSetMyCommands(botCommands...)
	_, err = bot.Request(setCommands)
	if err != nil {
		log.Fatal("Failed to set commands: ", err)
	}

	if cfg.ServerURL == "" {
		log.Fatal("Server URL is empty")
	}

	tgBot.ServerURL = cfg.ServerURL

	tgBot.BotAPI = bot

	return tgBot
}

func (bot *TgBot) SendMessage(chatId int64, message string) {
	const op = "tgbot.SendMessage"
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.DisableWebPagePreview = true
	_, err := bot.Send(msg)
	if err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
	}
}

func (bot *TgBot) GetCommand() []tgbotapi.BotCommand {
	const op = "tgbot.GetBotCommand"
	commands, err := bot.GetMyCommands()
	if err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
		return nil
	}
	return commands
}

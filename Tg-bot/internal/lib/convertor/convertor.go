package convertor

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-bot/internal/config"
)

func CommandConverter(commands []config.BotCommand) ([]tgbotapi.BotCommand, error) {
	var botCommands []tgbotapi.BotCommand
	for _, cmd := range commands {
		botCommands = append(botCommands, tgbotapi.BotCommand{Command: cmd.Command, Description: cmd.Description})
	}
	return botCommands, nil
}

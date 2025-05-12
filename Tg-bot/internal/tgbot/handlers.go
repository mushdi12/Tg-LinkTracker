package tgbot

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-bot/internal/lib/network"
	"tg-bot/internal/tgbot/commands"
	"tg-bot/internal/tgbot/user"
	"time"
)

func MainController(updates tgbotapi.UpdatesChannel, bot *TgBot, cancel <-chan struct{}, net *network.Network) {
	for {
		select {
		case update := <-updates:
			chatID := update.Message.Chat.ID
			username := update.Message.From.UserName
			if update.Message.IsCommand() {
				command := update.Message.Command()
				go commandHandler(username, command, chatID, bot, net)
			} else {
				message := update.Message.Text
				go messageHandler(message, chatID, bot, net)
			}
		case <-cancel:
			return
		}
	}
}

func commandHandler(username string, cmd string, chatID int64, bot *TgBot, net *network.Network) {
	op := "tgbot.handler.commandHandler"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	command, exists := commands.CommandRegistry[cmd]
	if !exists {
		bot.SendMessage(chatID, "Неизвестная команда! Воспользуйтесь /help")
	}

	commandCtx := commands.CommandContext{ChatId: chatID, Username: username}
	messageForUser := make(chan string)

	if cmd == "help" {
		commandCtx.BotCmd = bot.GetCommand()
	}

	go func() {
		msg := command.Execute(commandCtx, net)
		messageForUser <- msg
	}()

	select {
	case msg := <-messageForUser:
		bot.SendMessage(chatID, msg)
	case <-ctx.Done():
		_ = fmt.Errorf("%s: %w", op, errors.New("Команда "+cmd+" заняла слишком много времени!"))
		bot.SendMessage(chatID, "Произошла ошибка, попробуйте еще раз!")
	}

}

func messageHandler(message string, chatID int64, bot *TgBot, net *network.Network) {
	op := "tgbot.handler.messageHandler"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	state, err := user.GetState(chatID)

	if err != nil {
		bot.SendMessage(chatID, "Неизвестная команда! Воспользуйтесь /help")
		return
	}
	var command commands.Command

	if state < 2 {
		bot.SendMessage(chatID, "Неизвестная команда! Воспользуйтесь /help")
		return
	}

	if _, ok := user.TrackStates[state]; ok {
		command = &commands.TrackCommand{}
	} else if _, ok := user.UnTrackStates[state]; ok {
		command = &commands.UnTrackCommand{}
	} else {
		bot.SendMessage(chatID, "Неизвестная команда! Воспользуйтесь /help")
		return
	}

	messageForUser := make(chan string)

	go func() {
		msg := command.Execute(commands.CommandContext{ChatId: chatID, Message: message}, net)
		messageForUser <- msg
	}()

	select {
	case msg := <-messageForUser:
		bot.SendMessage(chatID, msg)
	case <-ctx.Done():
		_ = fmt.Errorf("%s: %w", op, errors.New("Обработка сообщения "+message+" заняла слишком много времени!"))
		bot.SendMessage(chatID, "Команда заняла слишком много времени, попробуйте еще раз!")
	}
}

package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	commands "tg-bot/internal/commands"
	user "tg-bot/internal/user"
	"time"
)

func MainController(updates tgbotapi.UpdatesChannel, bot *TgBot, stopChan <-chan struct{}) {
	for {
		select {
		case update := <-updates:
			chatID := update.Message.Chat.ID
			username := update.Message.From.UserName
			if update.Message.IsCommand() {
				command := update.Message.Command()
				go handleCommand(username, command, chatID, bot)
			} else {
				message := update.Message.Text
				go handleMessage(message, chatID, bot)
			}
		case <-stopChan:
			return
		}
	}
}

func handleCommand(username string, cmd string, chatID int64, bot *TgBot) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var command commands.Command
	var exists bool
	if command, exists = commands.CommandRegistry[cmd]; !exists {
		bot.SendMessage(chatID, "Неизвестная команда! Воспользуйтесь /help")
	}
	commandCtx := commands.CommandContext{ChatId: chatID, Username: username}
	messageForUser := make(chan string)

	if cmd == "help" {
		botCmds, err := bot.GetBotCommand()
		if err != nil {
			commandCtx.BotCmd = nil
		} else {
			commandCtx.BotCmd = botCmds
		}
	}

	go func() {
		msg := command.Execute(commandCtx)
		messageForUser <- msg
	}()

	select {
	case msg := <-messageForUser:
		bot.SendMessage(chatID, msg)
	case <-ctx.Done():
		bot.SendMessage(chatID, "Команда заняла слишком много времени, попробуйте еще раз!")
	}

}

func handleMessage(message string, chatID int64, bot *TgBot) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	state := user.GetState(chatID)
	_, existAddState := user.AddStates[state]
	_, existRemoveState := user.RemoveStates[state]
	if (existAddState || existRemoveState) && state > user.NONE {
		messageForUser := make(chan string)
		var command commands.Command
		if existAddState {
			command = &commands.TrackCommand{}
		} else {
			command = &commands.UnTrackCommand{}
		}
		go func() {
			msg := command.Execute(commands.CommandContext{ChatId: chatID, Message: message})
			messageForUser <- msg
		}()

		select {
		case msg := <-messageForUser:
			bot.SendMessage(chatID, msg)
		case <-ctx.Done():
			bot.SendMessage(chatID, "Команда заняла слишком много времени, попробуйте еще раз!")
		}
	} else {
		bot.SendMessage(chatID, "Неизвестная команда! Воспользуйтесь /help")
	}
}

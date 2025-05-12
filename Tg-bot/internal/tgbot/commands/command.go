package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"reflect"
	"tg-bot/internal/lib/network"
	. "tg-bot/internal/tgbot/user"
)

var CommandRegistry = map[string]Command{
	"start":   &StartCommand{},
	"help":    &HelpCommand{},
	"list":    &ListCommand{},
	"track":   &TrackCommand{},
	"untrack": &UnTrackCommand{},
}

type CommandContext struct {
	ChatId   int64
	Username string
	Message  string
	BotCmd   []tgbotapi.BotCommand
}

type Command interface {
	Execute(ctx CommandContext, net *network.Network) string
}

func setUserField(user *User, fieldName string, value string) {
	reflect.ValueOf(user).Elem().FieldByName(fieldName).SetString(value)
}

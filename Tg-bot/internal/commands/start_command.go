package commands

import (
	"sync"
	"tg-bot/internal/network"
	. "tg-bot/internal/user"
)

var (
	mu            sync.Mutex
	startMessages = map[int]string{
		400: "Ошибка со стороны клиента: попробуйте еще раз!",
		409: "Вы уже зарегистрированы!",
		500: "Ошибка со стороны сервера: попробуйте еще раз!",
		200: "Вы успешно авторизированы!"}
)

type StartCommand struct{}

func (cmd *StartCommand) Execute(ctx CommandContext) string {
	if !network.CheckUser(ctx.ChatId) {
		return listMessages[409]
	}

	code, err := network.AddUserRequest(ctx.ChatId, ctx.Username)
	if err != nil {
		return startMessages[400]
	}

	Users[ctx.ChatId] = &User{ChatId: ctx.ChatId, State: NONE}
	return startMessages[int(code)]
}

package commands

import (
	"tg-bot/internal/lib/network"
	"tg-bot/internal/lib/utils"
	"tg-bot/internal/tgbot/user"
)

type StartCommand struct{}

// todo: ошибка при добавление пользователся
func (cmd *StartCommand) Execute(ctx CommandContext, net *network.Network) string {
	if !utils.CheckDoubleCmmd(ctx.ChatId) {
		user.ResetState(ctx.ChatId)
		return "Нельзя во время другой команды использовать другие!"
	}
	answer := net.AddUserRequest(ctx.ChatId, ctx.Username)
	if answer != "Пользователь успешно добавлен" {
		return answer
	}
	user.Users[ctx.ChatId] = &user.User{ChatId: ctx.ChatId, State: user.NONE}
	return answer
}

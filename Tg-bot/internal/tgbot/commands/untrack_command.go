package commands

import (
	"tg-bot/internal/lib/network"
	"tg-bot/internal/tgbot/user"
)

type UnTrackCommand struct{}

func (cmd *UnTrackCommand) Execute(ctx CommandContext, net *network.Network) string {
	mu.Lock()
	defer mu.Unlock()

	if !net.CheckUserRequest(ctx.ChatId) {
		return "Такого пользователя не существует, сначала зарегистрируйтесь -> /start !"
	}

	userObj, _ := user.Users[ctx.ChatId]
	state := userObj.State

	if ctx.Message == "" && state != user.NONE {
		user.ResetState(ctx.ChatId)
		return "Ошибка! Действие команды отменено"
	}

	stmf := user.UnTrackStates[state]

	if stmf.FieldtoChange != "" {
		setUserField(userObj, stmf.FieldtoChange, ctx.Message)
	}

	userObj.State = stmf.NextState

	if state == user.WaitingUrlForRemove {
		return net.RemoveLinkRequest(ctx.ChatId, userObj.Link)
	}

	return stmf.Message
}

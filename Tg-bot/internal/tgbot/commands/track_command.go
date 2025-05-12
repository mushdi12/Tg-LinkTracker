package commands

import (
	"sync"
	"tg-bot/internal/lib/network"
	. "tg-bot/internal/tgbot/user"
)

var mu sync.Mutex

type TrackCommand struct{}

func (cmd *TrackCommand) Execute(ctx CommandContext, net *network.Network) string {
	mu.Lock()
	defer mu.Unlock()

	if !net.CheckUserRequest(ctx.ChatId) {
		return "Такого пользователя не существует, сначала зарегистрируйтесь -> /start !"
	}

	user, _ := Users[ctx.ChatId]
	state := user.State

	if ctx.Message == "" && state != NONE {
		ResetState(ctx.ChatId)
		return "Ошибка! Действие Команды отменено"
	}

	stmf := TrackStates[state]

	if stmf.FieldtoChange != "" {
		setUserField(user, stmf.FieldtoChange, ctx.Message)
	}

	user.State = stmf.NextState
	if state == WaitingHashtag {
		return net.AddLinkRequest(ctx.ChatId, user.Link, user.Category)
	}
	return stmf.Message
}

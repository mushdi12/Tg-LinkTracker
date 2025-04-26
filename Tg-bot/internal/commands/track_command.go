package commands

import (
	"tg-bot/internal/network"
	. "tg-bot/internal/user"
)

var (
	trackMessages = map[int]string{
		2:   "Ошибка! Действие Команды отменено",
		400: "Ошибка со стороны клиента: попробуйте еще раз!",
		409: "Такого пользователя не существует, сначала зарегистрируйтесь -> /start !",
		500: "Ошибка со стороны сервера: попробуйте еще раз!",
		200: "Ссылка успешно добавлена!"}
)

type TrackCommand struct{}

func (cmd *TrackCommand) Execute(ctx CommandContext) string {
	mu.Lock()
	defer mu.Unlock()

	if !network.CheckUser(ctx.ChatId) {
		return listMessages[409]
	}

	user, _ := Users[ctx.ChatId]
	state := user.State

	if ctx.Message == "" && state != NONE {
		ResetState(ctx.ChatId)
		return trackMessages[2]
	}

	stmf := AddStates[state]

	if stmf.FieldtoChange != "" {
		setUserField(user, stmf.FieldtoChange, ctx.Message)
	}

	user.State = stmf.NextState

	if state == WaitingHashtag {
		code, err := network.AddLinkRequest(ctx.ChatId, user.Link, user.Category)
		if err != nil {
			return startMessages[400]
		}
		return trackMessages[int(code)]
	}

	return stmf.Message
}

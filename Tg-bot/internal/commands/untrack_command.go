package commands

import (
	"tg-bot/internal/network"
	. "tg-bot/internal/user"
)

var (
	untrackMessages = map[int]string{
		1:   "Сначала зарегистрируйтесь -> /start",
		2:   "Ошибка! Действие команды отменено",
		400: "Ошибка со стороны клиента: попробуйте еще раз!",
		404: "У вас нет данной ссылки, проверьте через -> /list !",
		409: "Такого пользователя не существует, сначала зарегистрируйтесь -> /start !",
		500: "Ошибка со стороны сервера: попробуйте еще раз!",
		200: "Ссылка успешно удалена!"}
)

type UnTrackCommand struct{}

func (cmd *UnTrackCommand) Execute(ctx CommandContext) string {
	mu.Lock()
	defer mu.Unlock()

	if !network.CheckUser(ctx.ChatId) {
		return listMessages[409]
	}
	
	user, _ := Users[ctx.ChatId]
	state := user.State

	if ctx.Message == "" && state != NONE {
		ResetState(ctx.ChatId)
		return untrackMessages[2]
	}

	stmf := RemoveStates[state]

	if stmf.FieldtoChange != "" {
		setUserField(user, stmf.FieldtoChange, ctx.Message)
	}

	user.State = stmf.NextState

	if state == WaitingUrlForRemove {
		code, err := network.RemoveLinkRequest(ctx.ChatId, user.Link)
		if err != nil {
			return startMessages[400]
		}
		return trackMessages[int(code)]
	}

	return stmf.Message
}

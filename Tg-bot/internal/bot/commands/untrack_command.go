package commands

import (
	"log"
	"tg-bot/internal/network"
	. "tg-bot/internal/user"
)

type UnTrackCommand struct{}

func (cmd *UnTrackCommand) Execute(ctx CommandContext) string {
	mu.Lock()
	defer mu.Unlock()

	user, ok := Users[ctx.ChatId]
	if !ok {
		log.Printf("[UnTrackCommand.Execute] пользователь не найден в Users")
		return "Сначала введите /start"
	}

	state := user.State
	stmf := RemoveStates[state]

	if stmf.FieldtoChange != "" {
		setUserField(user, stmf.FieldtoChange, ctx.Message)
	}

	user.State = stmf.NextState

	if state == WaitingUrlForRemove {
		err := network.RemoveLinkRequest(ctx.ChatId, user.Link)
		if err != nil {
			log.Printf("[UnTrackCommand.Execute] ошибка при удалении ссылки: %v", err)
			return "Не удалось удалить ссылку, попробуйте снова"
		}
	}

	return stmf.Message
}

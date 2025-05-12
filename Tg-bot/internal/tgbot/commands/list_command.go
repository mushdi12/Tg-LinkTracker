package commands

import (
	"tg-bot/internal/lib/network"
	"tg-bot/internal/lib/utils"
	"tg-bot/internal/tgbot/user"
)

type ListCommand struct{}

func (cmd *ListCommand) Execute(ctx CommandContext, net *network.Network) string {
	if !utils.CheckDoubleCmmd(ctx.ChatId) {
		user.ResetState(ctx.ChatId)
		return "Нельзя во время другой команды использовать другие!"
	}
	return net.GetUsersLinkRequest(ctx.ChatId)
}

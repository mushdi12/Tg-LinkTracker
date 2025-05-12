package commands

import (
	"fmt"
	"strings"
	"tg-bot/internal/lib/network"
	"tg-bot/internal/lib/utils"
	"tg-bot/internal/tgbot/user"
)

type HelpCommand struct{}

func (cmd *HelpCommand) Execute(ctx CommandContext, _ *network.Network) string {
	commands := ctx.BotCmd
	if commands == nil {
		return "Произошла ошибка, попробуйте еще раз!"
	}

	if !utils.CheckDoubleCmmd(ctx.ChatId) {
		user.ResetState(ctx.ChatId)
		return "Нельзя во время другой команды использовать другие!"
	}

	var result strings.Builder
	for _, command := range commands {
		fmt.Fprintf(&result, "/%s - %s\n", command.Command, command.Description)
	}
	return result.String()
}

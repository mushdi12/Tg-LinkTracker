package commands

import (
	"fmt"
	"strings"
	"tg-bot/internal/network"
)

var (
	listMessages = map[int]string{
		2:   "Ошибка! Действие Команды отменено",
		3:   "У вас пока нет добавленных ссылок.",
		400: "Ошибка со стороны клиента: попробуйте еще раз!",
		409: "Такого пользователя не существует, сначала зарегистрируйтесь -> /start !",
		500: "Ошибка со стороны сервера: попробуйте еще раз!",
		200: "Ссылка успешно добавлена!"}
)

type ListCommand struct{}

func (cmd *ListCommand) Execute(ctx CommandContext) string {
	if !network.CheckUser(ctx.ChatId) {
		return listMessages[409]
	}
	userLinks, code, err := network.GetUsersLinkRequest(ctx.ChatId)
	if err != nil {
		return startMessages[400]
	}
	if code != 200 {
		return listMessages[int(code)]
	}
	return formatLinksByCategory(userLinks)
}

func formatLinksByCategory(links []network.Link) string {
	if len(links) == 0 {
		return listMessages[3]
	}

	categoryMap := make(map[string][]string)
	for _, link := range links {
		categoryMap[link.Category] = append(categoryMap[link.Category], link.URL)
	}

	var builder strings.Builder
	builder.WriteString("*Ваши ссылки*\n\n")
	for category, urls := range categoryMap {
		escapedCategory := escapeMarkdown(category)
		builder.WriteString(fmt.Sprintf("📂 *%s*\n", escapedCategory))
		for i, url := range urls {
			builder.WriteString(fmt.Sprintf("%d. [%s](%s)\n", i+1, url, url))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"*", "\\*", "_", "\\_", "[", "\\[", "]", "\\]",
		"(", "\\(", ")", "\\)",
		"~", "\\~", "`", "\\`", "#", "\\#",
	)
	return replacer.Replace(text)
}

package formatter

import (
	"fmt"
	"strings"
	"tg-bot/internal/lib/api"
)

func FormatLinksByCategory(links []api.Link) string {
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

package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"tg-bot/internal/lib/api"
	"tg-bot/internal/tgbot"
)

func HandleGitHub(bot *tgbot.TgBot) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctxGitRepo := ctx.Get("gitRepo")
		gitRepo, ok := ctxGitRepo.(api.GitRepo)

		if !ok {
			return ctx.String(http.StatusInternalServerError, "Не удалось получить ссылку из контекста")
		}

		// TODO поработать с msg
		bot.SendMessage(gitRepo.To, gitRepo.Msg)

		return ctx.String(http.StatusOK, "Ссылка добавлена")
	}
}

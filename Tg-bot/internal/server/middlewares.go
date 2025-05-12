package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"tg-bot/internal/lib/api"
)

func GitHubMiddleware() echo.MiddlewareFunc {
	op := "server.middlewares.GitHubMiddleware"
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var gitRepo api.GitRepo
			if err := ctx.Bind(&gitRepo); err != nil {
				log.Infof("%s: %w", op, err)
				return ctx.String(http.StatusBadRequest, "Невалидные данные пользователя")
			}

			if gitRepo.To == 0 {
				return ctx.String(http.StatusBadRequest, "Поле <to> обязательно")
			}

			if gitRepo.Name == "" {
				return ctx.String(http.StatusBadRequest, "Поле <Name> обязательно")
			}

			if gitRepo.Owner == "" {
				return ctx.String(http.StatusBadRequest, "Поле <Owner> обязательно")
			}

			if gitRepo.Msg == "" {
				return ctx.String(http.StatusBadRequest, "Поле <Msg> обязательно")
			}

			ctx.Set("gitRepo", gitRepo)

			return next(ctx)
		}
	}
}

package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"server-scrapper/internal/lib/api"
	"server-scrapper/internal/storage/postgres"
)

func ExtractUser(_ *postgres.Storage) echo.MiddlewareFunc {
	op := "server.middlewares.ExtractUser"
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var user api.User
			if err := ctx.Bind(&user); err != nil {
				log.Infof("%s: %w", op, err)
				return ctx.String(http.StatusBadRequest, "Невалидные данные пользователя")
			}

			if user.ChatID == 0 {
				return ctx.String(http.StatusBadRequest, "Поле <chat_id> обязателельно")
			}

			ctx.Set("user", user)

			return next(ctx)
		}
	}
}

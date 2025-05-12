package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"server-scrapper/internal/lib/api"
	"server-scrapper/internal/storage/postgres"
)

func ExtractUserAndCheck(storage *postgres.Storage) echo.MiddlewareFunc {
	op := "server.middlewares.ExtractUserAndCheck"
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var user api.User
			if err := ctx.Bind(&user); err != nil {
				log.Infof("%s: %w", op, err)
				return ctx.String(http.StatusBadRequest, "Невалидные данные пользователя")
			}

			if user.ChatID == 0 {
				return ctx.String(http.StatusBadRequest, "Поле <chat_id> обязателен")
			}

			exists, err := storage.UserExist(ctx.Request().Context(), user.ChatID)
			if err != nil {
				return ctx.String(http.StatusInternalServerError, "Ошибка при проверке пользователя")
			}
			if !exists {
				return ctx.String(http.StatusNotFound, "Пользователь не существует")
			}

			ctx.Set("user", user)

			return next(ctx)
		}
	}
}

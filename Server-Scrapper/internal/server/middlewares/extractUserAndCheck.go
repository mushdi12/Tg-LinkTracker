package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
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
				fmt.Errorf("%s: %w", op, err)
				return ctx.String(http.StatusBadRequest, "Невалидные данные пользователя")
			}

			if user.ChatID == 0 {
				return ctx.String(http.StatusBadRequest, "Поле <chat_id> обязателен")
			}

			if storage.UserExist(ctx.Request().Context(), user.ChatID) {
				return ctx.String(http.StatusConflict, "Такой пользователь уже существует")
			}

			ctx.Set("user", user)

			return next(ctx)
		}
	}
}

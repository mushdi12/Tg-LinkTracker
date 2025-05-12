package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"server-scrapper/internal/storage/postgres"
	"strconv"
)

func ExtractChatIdAndChecke(storage *postgres.Storage) echo.MiddlewareFunc {
	op := "server.middlewares.ExtractChatIdAndChecke"
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			chatIDStr := ctx.QueryParam("chatId")
			chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
			if err != nil {
				log.Infof("%s: %w", op, err)
				return ctx.String(http.StatusBadRequest, "Некорректный chatId")
			}

			exists, err := storage.UserExist(ctx.Request().Context(), chatID)
			if err != nil {
				return ctx.String(http.StatusInternalServerError, "Ошибка при проверке пользователя")
			}
			if !exists {
				return ctx.String(http.StatusNotFound, "Пользователь не существует")
			}

			ctx.Set("chatId", chatID)

			return next(ctx)
		}
	}
}

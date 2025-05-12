package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"server-scrapper/internal/lib/api"
	"server-scrapper/internal/storage/postgres"
)

func ExtractLink(storage *postgres.Storage) echo.MiddlewareFunc {
	op := "server.middlewares.ExtractLink"
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var linkRequest api.LinkRequest
			if err := ctx.Bind(&linkRequest); err != nil {
				log.Infof("%s: %w", op, err)
				return ctx.String(http.StatusBadRequest, "Невалидные данные пользователя")
			}

			if linkRequest.ChatID == 0 {
				return ctx.String(http.StatusBadRequest, "Поле <chat_id> обязательно")
			}

			if linkRequest.Link == "" {
				return ctx.String(http.StatusBadRequest, "Поле <link> обязательно")
			}

			exists, err := storage.UserExist(ctx.Request().Context(), linkRequest.ChatID)
			if err != nil {
				return ctx.String(http.StatusInternalServerError, "Ошибка при проверке пользователя")
			}
			if !exists {
				return ctx.String(http.StatusNotFound, "Пользователь не существует")
			}

			ctx.Set("linkRequest", linkRequest)

			return next(ctx)
		}
	}
}

package removeLink

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"server-scrapper/internal/storage"
	"strconv"
)

type LinkDeleter interface {
	DeleteLink(ctx context.Context, chatID int64, linkUrl string) error
}

// TODO: ДОБАВИТЬ СВОИ КОНТЕКСТЫ И ТАЙМАУТЫ И ГОРУТИНЫ
func Handler(linkDeleter LinkDeleter) echo.HandlerFunc {
	op := "server.handlers.removeLink.handler"
	return func(ctx echo.Context) error {
		chatIDStr := ctx.QueryParam("chatId")
		userLink := ctx.QueryParam("link")
		chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)

		err := linkDeleter.DeleteLink(ctx.Request().Context(), chatID, userLink)
		if err != nil {
			if errors.Is(err, storage.ErrLinkNotFound) {
				return ctx.String(http.StatusNotFound, "У пользователя нет такой ссылки")
			}
			log.Info("%s: %w", op, err)
			return ctx.String(http.StatusInternalServerError, "Ошибка при удалении ссылки")
		}
		return ctx.String(http.StatusOK, "Ссылка удалена")
	}
}

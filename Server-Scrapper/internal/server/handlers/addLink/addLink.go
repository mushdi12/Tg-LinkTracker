package addLink

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"server-scrapper/internal/lib/api"
)

type LinkInserter interface {
	InsertLink(ctx context.Context, chatID int64, linkUrl string, linkCategory string) bool
}

func Handler(linkInserter LinkInserter) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctxLinkRequest := ctx.Get("linkRequest")
		linkRequest, ok := ctxLinkRequest.(api.LinkRequest)
		if !ok {
			return ctx.String(http.StatusInternalServerError, "Не удалось получить ссылку из контекста")
		}

		if !linkInserter.InsertLink(ctx.Request().Context(), linkRequest.ChatID, linkRequest.Link, linkRequest.Category) {
			return ctx.String(http.StatusInternalServerError, "Ошибка при добавлении ссылки")
		}
		return ctx.String(http.StatusOK, "Ссылка добавлена")
	}
}

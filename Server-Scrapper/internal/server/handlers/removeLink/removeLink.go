package removeLink

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type LinkDeleter interface {
	DeleteLink(ctx context.Context, chatID int64, linkUrl string) (bool, string)
}

func Handler(linkDeleter LinkDeleter) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		chatIDStr := ctx.QueryParam("chatId")
		userLink := ctx.QueryParam("link")
		chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)

		ok, answer := linkDeleter.DeleteLink(ctx.Request().Context(), chatID, userLink)
		if !ok {
			return ctx.String(http.StatusInternalServerError, answer)
		}

		return ctx.String(http.StatusOK, "Ссылка удалена")
	}
}

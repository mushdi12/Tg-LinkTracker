package getLinks

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"server-scrapper/internal/lib/api"
	"server-scrapper/internal/lib/convertor"
	"strconv"
)

type LinkGetter interface {
	GetUserLinks(ctx context.Context, chatID int64) (map[string]string, error)
}

// TODO: ДОБАВИТЬ СВОИ КОНТЕКСТЫ И ТАЙМАУТЫ И ГОРУТИНЫ
func Handler(linkGetter LinkGetter) echo.HandlerFunc {
	op := "server.handlers.getLink.handler"
	return func(ctx echo.Context) error {
		chatIDStr := ctx.QueryParam("chatId")
		chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)

		userLinks, err := linkGetter.GetUserLinks(ctx.Request().Context(), chatID)

		if err != nil {
			log.Info("%s: %w", op, err)
			return ctx.String(http.StatusInternalServerError, "Ошибка при получении ссылок")
		}

		if len(userLinks) == 0 {
			return ctx.String(http.StatusNotFound, "Ссылки не найдены для данного пользователя")
		}

		links := convertor.ConvertLinks(userLinks)

		resp := api.LinksResponse{
			Links: links,
		}

		return ctx.JSON(http.StatusOK, resp)
	}
}

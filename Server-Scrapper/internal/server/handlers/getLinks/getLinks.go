package getLinks

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"server-scrapper/internal/lib/api"
	"server-scrapper/internal/storage/postgres"
	"strconv"
)

type LinkGetter interface {
	GetUserLinks(ctx context.Context, chatID int64) UserLinks
}

func Handler(storage *postgres.Storage) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		chatIDStr := ctx.QueryParam("chatId")
		chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)

		userLinks := storage.GetUserLinks(ctx.Request().Context(), chatID)

		var links []api.Link
		for _, l := range userLinks {
			links = append(links, api.Link{
				URL:      l.LinkUrl,
				Category: l.LinkCategory,
			})
		}

		resp := api.LinksResponse{
			Links: links,
		}

		return ctx.JSON(http.StatusOK, resp)
	}
}

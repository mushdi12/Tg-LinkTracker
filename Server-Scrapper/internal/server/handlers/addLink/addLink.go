package addLink

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"server-scrapper/internal/lib/api"
	"server-scrapper/internal/lib/github"
	"server-scrapper/internal/storage"
	"time"
)

type LinkInserter interface {
	InsertLink(ctx context.Context, chatID int64, linkUrl string, linkCategory string) error
	InsertRepo(ctx context.Context, chatID int64, repoOwner, repoName, lastSHA string, lastIssue, lastPR int) error
}

// TODO: ДОБАВИТЬ СВОИ КОНТЕКСТЫ И ТАЙМАУТЫ И ГОРУТИНЫ
func Handler(linkInserter LinkInserter, gitClient *github.GitHub) echo.HandlerFunc {
	op := "server.handlers.addLink.handler"
	return func(ctx echo.Context) error {
		ctxLinkRequest := ctx.Get("linkRequest")
		linkRequest, ok := ctxLinkRequest.(api.LinkRequest)

		if !ok {
			return ctx.String(http.StatusInternalServerError, "Не удалось получить ссылку из контекста")
		}

		err := linkInserter.InsertLink(ctx.Request().Context(), linkRequest.ChatID, linkRequest.Link, linkRequest.Category)

		if err != nil {
			if errors.Is(err, storage.ErrLinkExists) {
				log.Infof("%s: %w", op, err)
				return ctx.String(http.StatusConflict, "Такая ссылка уже существует")
			}
			log.Infof("%s: %w", op, err)
			return ctx.String(http.StatusInternalServerError, "Ошибка при добавлении ссылки")
		}

		gitCtx, cancel := context.WithTimeout(ctx.Request().Context(), 10*time.Second)
		defer cancel()

		info, err := gitClient.GetFirstRepoInfo(gitCtx, linkRequest.Link, 10*time.Second)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = linkInserter.InsertRepo(ctx.Request().Context(), linkRequest.ChatID,
			info.RepoOwner, info.RepoName, info.LastSHA, info.LastIssue, info.LastPR)

		if err != nil {
			if errors.Is(err, storage.ErrLinkExists) {
				log.Infof("%s: %w", op, err)
				return ctx.String(http.StatusConflict, "Такая ссылка уже существует")
			}
			log.Infof("%s: %w", op, err)
			return ctx.String(http.StatusInternalServerError, "Ошибка при добавлении ссылки")
		}

		return ctx.String(http.StatusOK, "Ссылка добавлена")
	}
}

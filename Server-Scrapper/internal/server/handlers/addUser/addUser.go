package addUser

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"server-scrapper/internal/lib/api"
	"server-scrapper/internal/storage"
)

type UserInserter interface {
	InsertUser(ctx context.Context, chatID int64, username string) error
}

// TODO: ДОБАВИТЬ СВОИ КОНТЕКСТЫ И ТАЙМАУТЫ И ГОРУТИНЫ
func Handler(userInserter UserInserter) echo.HandlerFunc {
	op := "server.handlers.addUser.handler"
	return func(ctx echo.Context) error {
		ctxUser := ctx.Get("user")
		user, ok := ctxUser.(api.User)
		if !ok {
			return ctx.String(http.StatusInternalServerError, "Не удалось получить пользователя из контекста")
		}

		err := userInserter.InsertUser(ctx.Request().Context(), user.ChatID, user.Username)
		if err != nil {
			if errors.Is(err, storage.ErrUserExists) {
				return ctx.String(http.StatusConflict, "Пользователь уже существует")
			}
			log.Infof("%s: %s", op, err)
			return ctx.String(http.StatusInternalServerError, "Ошибка при добавлении пользователя")
		}
		return ctx.String(http.StatusOK, "Пользователь добавлен")
	}
}

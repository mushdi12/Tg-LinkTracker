package addUser

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"server-scrapper/internal/lib/api"
)

type UserInserter interface {
	InsertUser(ctx context.Context, chatID int64, username string) bool
}

func Handler(userInserter UserInserter) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctxUser := ctx.Get("user")
		user, ok := ctxUser.(api.User)
		if !ok {
			return ctx.String(http.StatusInternalServerError, "Не удалось получить пользователя из контекста")
		}

		if !userInserter.InsertUser(ctx.Request().Context(), user.ChatID, user.Username) {
			return ctx.String(http.StatusInternalServerError, "Ошибка при добавлении пользователя")
		}
		return ctx.String(http.StatusOK, "Пользователь добавлен")
	}
}

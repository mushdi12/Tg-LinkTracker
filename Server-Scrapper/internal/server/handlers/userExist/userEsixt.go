package userExist

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"server-scrapper/internal/lib/api"
)

type UserExister interface {
	UserExist(ctx context.Context, chatID int64) bool
}

func Handler(userExister UserExister) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctxUser := ctx.Get("user")
		user, ok := ctxUser.(api.User)
		if !ok {
			return ctx.String(http.StatusInternalServerError, "Не удалось получить пользователя из контекста")
		}

		if !userExister.UserExist(ctx.Request().Context(), user.ChatID) {
			return ctx.String(http.StatusNotFound, "Пользователь не существует")
		}
		return ctx.String(http.StatusOK, "Пользователь существует")
	}
}

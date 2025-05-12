package userExist

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"server-scrapper/internal/lib/api"
)

type UserExister interface {
	UserExist(ctx context.Context, chatID int64) (bool, error)
}

// TODO: ДОБАВИТЬ СВОИ КОНТЕКСТЫ И ТАЙМАУТЫ И ГОРУТИНЫ
func Handler(userExister UserExister) echo.HandlerFunc {
	op := "server.handlers.userExist.handler"
	return func(ctx echo.Context) error {
		ctxUser := ctx.Get("user")
		user, ok := ctxUser.(api.User)
		if !ok {
			return ctx.String(http.StatusInternalServerError, "Не удалось получить пользователя из контекста")
		}
		exists, err := userExister.UserExist(ctx.Request().Context(), user.ChatID)
		if err != nil {
			log.Info("%s: %w", op, err)
			return ctx.String(http.StatusInternalServerError, "Ошибка при проверке пользователя")
		}
		if !exists {
			return ctx.String(http.StatusNotFound, "Пользователь не существует")
		}
		return ctx.String(http.StatusOK, "Пользователь существует")
	}
}

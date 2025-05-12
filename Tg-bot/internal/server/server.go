package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"tg-bot/internal/tgbot"
)

type Server struct {
	*echo.Echo
	Port string
}

func New(bot *tgbot.TgBot) *Server {
	e := echo.New()

	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//TODO: ДОБАВИТЬ REDDIT И ТД
	e.POST("githubInfo", HandleGitHub(bot), GitHubMiddleware())

	return &Server{Echo: e, Port: ":8081"}
}

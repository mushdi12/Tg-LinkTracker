package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"server-scrapper/internal/lib/github"
	"server-scrapper/internal/server/handlers/addLink"
	"server-scrapper/internal/server/handlers/addUser"
	"server-scrapper/internal/server/handlers/getLinks"
	"server-scrapper/internal/server/handlers/removeLink"
	"server-scrapper/internal/server/handlers/userExist"
	"server-scrapper/internal/server/middlewares"
	"server-scrapper/internal/storage/postgres"
)

type Server struct {
	*echo.Echo
}

func New(storage *postgres.Storage, gitClient *github.GitHub) *Server {
	e := echo.New()

	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("user/addUser", addUser.Handler(storage), middlewares.ExtractUserAndCheck(storage))
	e.POST("user/isExist", userExist.Handler(storage), middlewares.ExtractUser(storage))
	e.POST("user/addLink", addLink.Handler(storage, gitClient), middlewares.ExtractLink(storage))
	e.GET("user/getLinks", getLinks.Handler(storage), middlewares.ExtractChatIdAndChecke(storage))
	e.DELETE("user/removeLink", removeLink.Handler(storage), middlewares.ExtractChatIdLinkAndCheck(storage))

	return &Server{Echo: e}
}

package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server-scrapper/internal/config"
	"server-scrapper/internal/lib/github"
	"server-scrapper/internal/lib/telegram"
	"server-scrapper/internal/server"
	"server-scrapper/internal/storage/postgres"
	"server-scrapper/internal/watcher"
	"syscall"
	"time"
)

type App struct {
	Storage *postgres.Storage
	Watcher *watcher.Watcher
	Server  *server.Server
}

func New(cfg *config.Config) *App {

	storage := postgres.New(cfg.StorageData)
	tgClient := telegram.New(cfg.TelegramBot)
	gitClient := github.New()
	watcher := watcher.New(storage, tgClient, gitClient, cfg.Watcher)
	server := server.New(storage, gitClient)

	return &App{Storage: storage, Watcher: watcher, Server: server}
}

func (app *App) Start() {

	app.Watcher.StartMonitoring()

	go func() {
		if err := app.Server.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Printf("ошибка при завершении сервера", "error", err)
	}
}

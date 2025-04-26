package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"server-scrapper/internal/config"
	"server-scrapper/internal/lib/logger/sl"
	"server-scrapper/internal/server"
	"server-scrapper/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad("configs/bd_config.json")

	log := setupLogger()
	_ = log
	storage, err := postgres.New(cfg.StorageData)
	if err != nil {
		log.Error("failed to connect to storage", sl.Err(err))
	}

	router := server.New(storage)
	if err := router.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}

func setupLogger() *slog.Logger {
	var logger *slog.Logger
	// задел под будущее, но пока все будет выводиться в консоль текстом, с уровнем - дебаг
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}

// Package main starts the URL shortener application.
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"shortener/internal/app"
	"shortener/internal/config"
	"shortener/internal/db"
	"shortener/internal/httpserver"
	"shortener/internal/services/shortener"
	"syscall"
)

func main() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.New(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	service := shortener.NewService(database.Links, database.LinkVisits, &cfg.App)
	server := httpserver.New(service, &cfg.HTTP)
	app := app.New(server, database, cfg)

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

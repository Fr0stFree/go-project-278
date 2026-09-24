// Package main starts the URL shortener application.
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"shortener/internal/app"
	"shortener/internal/config"
	"shortener/internal/db"
	"shortener/internal/httpserver"
	"shortener/internal/integrations/sentry"
	"shortener/internal/services/shortener"
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

	service := shortener.NewService(database.Links, database.LinkVisits)
	server := httpserver.New(service, database, &cfg.HTTP, cfg.App.BaseURL)

	if err := sentry.Connect(cfg.HTTP.Sentry); err != nil {
		log.Fatal(err)
	}

	app := app.New(server, database, cfg)

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

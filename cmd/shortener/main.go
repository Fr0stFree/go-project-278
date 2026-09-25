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

	sentryIntegration, err := sentry.New(cfg.HTTP.Sentry)
	if err != nil {
		log.Fatal(err)
	}

	server := httpserver.New(service, database, &cfg.HTTP, sentryIntegration.Middleware())

	app := app.New(server, &cfg.App)
	app.AddShutdownCallback(func() error {
		sentryIntegration.Flush()

		return nil
	})
	app.AddShutdownCallback(database.Close)

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

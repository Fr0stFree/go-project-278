// Package main starts the URL shortener application.
package main

import (
	"context"
	"log"
	"os/signal"
	"shortener/internal/app"
	"shortener/internal/config"
	"shortener/internal/db"
	"shortener/internal/httpserver"
	"shortener/internal/services/shortener"
	"syscall"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.New(&cfg.DataBase)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	service := shortener.NewService(database.Link, database.LinkVisit, &cfg.App)
	server := httpserver.New(service, &cfg.HTTP)
	app := app.New(server, database, service, cfg)

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

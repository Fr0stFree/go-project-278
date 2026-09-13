// Package main starts the URL shortener application.
package main

import (
	"log"
	"shortener/internal/app"
	"shortener/internal/config"
	"shortener/internal/db"
	"shortener/internal/httpserver"
	"shortener/internal/services/shortener"
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

	service := shortener.NewService(database.Link, database.LinkVisit, &cfg.App)
	server := httpserver.New(service, &cfg.HTTP)
	runner := app.New(server)

	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}
}

// Package main starts the URL shortener application.
package main

import (
	"log"
	"shortener/internal/app"
	"shortener/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	runner, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}
}

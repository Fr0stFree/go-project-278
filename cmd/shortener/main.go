// Package main is the entry point of the URL shortener application.
package main

import (
	"log"
	"shortener/internal/app"
	"shortener/internal/config"
)

func main() {
	cfg := config.New()

	runner, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}
}

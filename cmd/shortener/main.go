package main

import (
	"log"
	"shortener/internal/app"
	"shortener/internal/config"
)

func main() {
	cfg := config.New()
	runner := app.New(cfg)

	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}
}

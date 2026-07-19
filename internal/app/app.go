// Package app provides the main application structure and methods to run it
package app

import (
	"errors"
	"fmt"
	"net/http"
	"shortener/internal/config"
	"shortener/internal/httpserver"
	"shortener/internal/shortener"
	"shortener/internal/storage/memory"
)

// App represents the main application structure
type App struct {
	server *http.Server
	cfg    *config.Config
}

// New creates a new App instance with the provided configuration
func New(cfg *config.Config) *App {
	repo := memory.NewRepository()
	shortener := shortener.NewService(repo, cfg.BaseURL)
	server := httpserver.New(cfg.HTTP, shortener)

	return &App{
		server: server,
		cfg:    cfg,
	}
}

// Run starts the HTTP server and listens for incoming requests
func (a *App) Run() error {
	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run HTTP server: %w", err)
	}

	return nil
}

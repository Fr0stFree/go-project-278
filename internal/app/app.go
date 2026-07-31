// Package app provides the main application structure and methods to run it
package app

import (
	"errors"
	"fmt"
	"net/http"
	"shortener/internal/config"
)

// App represents the main application structure
type App struct {
	server *http.Server
	cfg    *config.Root
}

// New creates a new App instance with the provided configuration
func New(cfg *config.Root) *App {
	repos := prepareRepos(&cfg.DataBase)
	services :=  prepareServices(&cfg.App, repos)
	server := prepareServer(&cfg.HTTP, services)
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

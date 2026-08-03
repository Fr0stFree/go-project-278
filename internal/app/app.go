// Package app wires application dependencies and runs the HTTP server.
package app

import (
	"errors"
	"fmt"
	"net/http"
	"shortener/internal/config"
)

// App owns the configured HTTP server.
type App struct {
	server *http.Server
	cfg    *config.Root
}

// New builds repositories, services, and the HTTP server from the provided configuration.
func New(cfg *config.Root) (*App, error) {
	repos, err := buildRepos(&cfg.DataBase)
	if err != nil {
		return nil, fmt.Errorf("build repositories: %w", err)
	}

	services := buildServices(&cfg.App, repos)
	server := buildServer(&cfg.HTTP, services)

	return &App{
		server: server,
		cfg:    cfg,
	}, nil
}

// Run starts the HTTP server and returns unexpected server errors.
func (a *App) Run() error {
	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run HTTP server: %w", err)
	}

	return nil
}

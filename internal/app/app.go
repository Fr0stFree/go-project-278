package app

import (
	"errors"
	"fmt"
	"net/http"
	"shortener/internal/config"
	"shortener/internal/httpserver"
)

type App struct {
	server *http.Server
	cfg    *config.Config
}

func New(cfg *config.Config) *App {
	return &App{
		server: httpserver.New(cfg.HTTP),
	}
}

func (a *App) Run() error {
	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run HTTP server: %w", err)
	}
	return nil
}

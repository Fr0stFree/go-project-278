// Package app wires application dependencies and runs the HTTP server.
package app

import (
	"errors"
	"fmt"
	"net/http"
)

// App owns the configured HTTP server.
type App struct {
	server *http.Server
}

// New builds repositories, services, and the HTTP server from the provided configuration.
func New(server *http.Server) *App {
	return &App{server: server}
}

// Run starts the HTTP server and returns unexpected server errors.
func (a *App) Run() error {
	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run HTTP server: %w", err)
	}

	return nil
}

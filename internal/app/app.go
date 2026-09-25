// Package app wires application dependencies and runs the HTTP server.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"shortener/internal/config"
)

type server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

// App owns the configured HTTP server.
type App struct {
	server            server
	cfg               *config.App
	shutdownCallbacks []func() error
}

// New builds an application with the given server and lifecycle configuration.
func New(server server, cfg *config.App) *App {
	return &App{server: server, cfg: cfg}
}

// AddShutdownCallback registers a callback to be executed during application shutdown.
func (a *App) AddShutdownCallback(callback func() error) {
	a.shutdownCallbacks = append(a.shutdownCallbacks, callback)
}

// Run starts the HTTP server and returns unexpected server errors.
func (a *App) Run(ctx context.Context) error {
	var errs []error

	errCh := make(chan error, 1)

	go func() {
		slog.Info("App started successfully")

		errCh <- a.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
		defer cancel()

		err := a.server.Shutdown(shutdownCtx)
		if err != nil {
			errs = append(errs, fmt.Errorf("shutdown HTTP server: %w", err))
		}
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs = append(errs, fmt.Errorf("run HTTP server: %w", err))
		}
	}

	for _, callback := range a.shutdownCallbacks {
		if err := callback(); err != nil {
			errs = append(errs, fmt.Errorf("on shutdown: %w", err))
		}
	}

	return errors.Join(errs...)
}

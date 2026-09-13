// Package app wires application dependencies and runs the HTTP server.
package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"shortener/internal/config"
	"shortener/internal/db"
	"shortener/internal/services/shortener"
)

// App owns the configured HTTP server.
type App struct {
	server    *http.Server
	shortener *shortener.Service
	db        *db.Database
	config    *config.Root
}

// New builds repositories, services, and the HTTP server from the provided configuration.
func New(server *http.Server, db *db.Database, shortener *shortener.Service, cfg *config.Root) *App {
	return &App{server: server, db: db, shortener: shortener, config: cfg}
}

// Run starts the HTTP server and returns unexpected server errors.
func (a *App) Run(ctx context.Context) error {
	var errs []error

	errCh := a.startHTTPServer()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.HTTP.ShutdownTimeout)
		defer cancel()

		err := a.server.Shutdown(shutdownCtx)
		if err != nil {
			errs = append(errs, fmt.Errorf("shutdown HTTP server: %w", err))
		}

		err = a.db.Close()
		if err != nil {
			errs = append(errs, fmt.Errorf("close database: %w", err))
		}

		return errors.Join(errs...)
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs = append(errs, fmt.Errorf("run HTTP server: %w", err))
		}

		err = a.db.Close()
		if err != nil {
			errs = append(errs, fmt.Errorf("close database: %w", err))
		}

		return errors.Join(errs...)
	}
}

func (a *App) startHTTPServer() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.server.ListenAndServe()
	}()

	return errCh
}

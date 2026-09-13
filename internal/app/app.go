// Package app wires application dependencies and runs the HTTP server.
package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"shortener/internal/config"
)

type server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

type database interface {
	Close() error
}

// App owns the configured HTTP server.
type App struct {
	server server
	db     database
	cfg    *config.Root
}

// New builds App with the given server and database, applying options.
func New(server server, db database, cfg *config.Root) *App {
	return &App{server: server, db: db, cfg: cfg}
}

// Run starts the HTTP server and returns unexpected server errors.
func (a *App) Run(ctx context.Context) error {
	var errs []error

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.HTTP.ShutdownTimeout)
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

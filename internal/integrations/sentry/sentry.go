// Package sentry provides integration with the Sentry error tracking service.
package sentry

import (
	"fmt"
	"log/slog"

	"github.com/getsentry/sentry-go"

	"shortener/internal/config"
)

// Connect initializes the Sentry client with the provided configuration.
func Connect(cfg config.Sentry) error {
	if !cfg.IsEnabled {
		slog.Info("Sentry is disabled")

		return nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.DSN,
		Environment: cfg.Environment,
	})
	if err != nil {
		return fmt.Errorf("connect to Sentry: %w", err)
	}

	slog.Info("Sentry connected successfully")

	return nil
}

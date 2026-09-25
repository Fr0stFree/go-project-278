// Package sentry provides integration with the Sentry error tracking service.
package sentry

import (
	"fmt"
	"log/slog"

	sdk "github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"

	"shortener/internal/config"
)

// Integration defines the interface for Sentry integration in the application.
type Integration interface {
	// CaptureException reports a handled request error through its request-scoped Sentry hub.
	CaptureException(ctx *gin.Context, err error)

	// Flush waits for queued events to be delivered before the application exits.
	Flush()
}

// New initializes a Sentry integration based on the provided configuration.
func New(cfg config.Sentry) (Integration, error) {
	if !cfg.IsEnabled {
		return &noopIntegration{}, nil
	}

	err := sdk.Init(sdk.ClientOptions{
		Dsn:         cfg.DSN,
		Environment: cfg.Environment,
	})
	if err != nil {
		return nil, fmt.Errorf("connect to Sentry: %w", err)
	}

	slog.Info("Sentry connected successfully")

	return &activeIntegration{config: cfg}, nil
}

type activeIntegration struct {
	config config.Sentry
}

func (a *activeIntegration) CaptureException(ctx *gin.Context, err error) {
	hub := sentrygin.GetHubFromContext(ctx)
	if hub == nil {
		hub = sdk.CurrentHub()
	}

	hub.CaptureException(err)
}

func (a *activeIntegration) Flush() {
	if !sdk.Flush(a.config.FlushTimeout) {
		slog.Warn("Sentry flush timed out", slog.Duration("timeout", a.config.FlushTimeout))
	}
}

type noopIntegration struct{}

func (d *noopIntegration) CaptureException(_ *gin.Context, _ error) {}
func (d *noopIntegration) Flush()                                   {}

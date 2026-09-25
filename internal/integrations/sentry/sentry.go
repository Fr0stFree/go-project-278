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
	// Flush waits for queued events to be delivered before the application exits.
	Flush()

	// Middleware returns a Gin middleware that integrates Sentry error tracking.
	Middleware() gin.HandlerFunc
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

func (a *activeIntegration) captureException(ctx *gin.Context, err error) {
	hub := sentrygin.GetHubFromContext(ctx)
	if hub == nil {
		hub = sdk.CurrentHub()
	}

	hub.CaptureException(err)
}

func (a *activeIntegration) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
			Timeout:         a.config.FlushTimeout,
		})(ctx)

		for _, err := range ctx.Errors {
			a.captureException(ctx, err.Err)
		}
	}
}

func (a *activeIntegration) Flush() {
	if !sdk.Flush(a.config.FlushTimeout) {
		slog.Warn("Sentry flush timed out", slog.Duration("timeout", a.config.FlushTimeout))
	}
}

type noopIntegration struct{}

func (*noopIntegration) Flush() {}
func (*noopIntegration) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

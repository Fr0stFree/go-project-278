package sentry

import (
	"shortener/internal/config"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// Middleware returns a Gin middleware that integrates Sentry error tracking.
func Middleware(cfg config.Sentry, integration Integration) gin.HandlerFunc {
	if !cfg.IsEnabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	sentryMiddleware := sentrygin.New(sentrygin.Options{
		Repanic:         true,
		WaitForDelivery: false,
		Timeout:         cfg.FlushTimeout,
	})

	return func(c *gin.Context) {
		sentryMiddleware(c)

		for _, err := range c.Errors {
			integration.CaptureException(c, err.Err)
		}
	}
}

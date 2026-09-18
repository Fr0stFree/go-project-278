package sentry

import (
	"shortener/internal/config"

	sentry "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// Middleware returns a Gin middleware that integrates Sentry error tracking.
func Middleware(cfg config.Sentry) gin.HandlerFunc {
	if !cfg.IsEnabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return sentry.New(sentry.Options{
		Repanic:         true,
		WaitForDelivery: false,
		Timeout:         cfg.FlushTimeout,
	})
}

// Package health exposes health check handlers.
package health

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ReadinessChecker verifies that application dependencies are available.
type ReadinessChecker interface {
	PingContext(ctx context.Context) error
}

type handler struct {
	checker            ReadinessChecker
	healthcheckTimeout time.Duration
}

func (h *handler) ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

func (h *handler) health(ctx *gin.Context) {
	checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), h.healthcheckTimeout)
	defer cancel()

	if err := h.checker.PingContext(checkCtx); err != nil {
		slog.ErrorContext(ctx, "health check failed", slog.Any("error", err))
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"ok": false})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

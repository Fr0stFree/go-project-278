// Package health provides HTTP handlers for health check endpoints.
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

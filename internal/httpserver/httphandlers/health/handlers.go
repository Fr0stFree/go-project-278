// Package health exposes health check handlers.
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

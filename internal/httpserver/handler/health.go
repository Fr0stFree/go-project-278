package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthHandler struct{}

func newHealthHandler() *healthHandler {
	return &healthHandler{}
}

func (h *healthHandler) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

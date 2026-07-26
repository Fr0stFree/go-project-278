package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

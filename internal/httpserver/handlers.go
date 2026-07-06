package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func pingHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

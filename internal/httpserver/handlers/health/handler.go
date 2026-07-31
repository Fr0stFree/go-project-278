package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func RegisterRoutes(router gin.IRouter) {
	router.GET("/ping", ping)
}

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

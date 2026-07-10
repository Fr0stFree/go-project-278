package httpserver

import (
	"github.com/gin-gonic/gin"

	"shortener/internal/httpserver/handler"
)

func newRouter(h *handler.Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/ping", h.Health.Ping)

	return router
}

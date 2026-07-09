package httpserver

import (
	"github.com/gin-gonic/gin"

	"shortener/internal/httpserver/handler"
)

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	healthHandler := handler.NewHealthHandler()
	router.GET("/ping", healthHandler.Ping)

	return router
}

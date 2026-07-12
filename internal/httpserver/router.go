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
	router.POST("/api/links", h.Link.Create)
	router.GET("/api/links/:id", h.Link.Get)
	router.GET("/api/links", h.Link.List)
	// router.PUT("/api/links/:id", h.Link.Update)
	// router.DELETE("/api/links/:id", h.Link.Delete)

	return router
}

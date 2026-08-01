package link

import (
	"shortener/internal/services/metrics"
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the link-related routes with the provided router.
func RegisterRoutes(shortener *shortener.Service, metrics *metrics.Service, router gin.IRouter) {
	h := &handler{shortener: shortener, metrics: metrics}
	router.POST("/api/links", h.create)
	router.GET("/api/links/:id", h.get)
	router.GET("/api/links", h.list)
	router.PUT("/api/links/:id", h.update)
	router.DELETE("/api/links/:id", h.delete)
	router.GET("/r/:short_name", h.redirect)
}

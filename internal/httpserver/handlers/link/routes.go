package link

import (
	"github.com/gin-gonic/gin"
	"shortener/internal/services/shortener"
)

// RegisterRoutes registers the link-related routes with the provided router.
func RegisterRoutes(service *shortener.Service, router gin.IRouter) {
	h := &handler{Service: service}
	router.POST("/api/links", h.create)
	router.GET("/api/links/:id", h.get)
	router.GET("/api/links", h.list)
	router.PUT("/api/links/:id", h.update)
	router.DELETE("/api/links/:id", h.delete)
}

package link

import (
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts link management and redirect routes.
func RegisterRoutes(shortener *shortener.Service, router gin.IRouter) {
	h := &handler{shortener: shortener}
	router.POST("/api/links", h.create)
	router.GET("/api/links/:id", h.get)
	router.GET("/api/links", h.list)
	router.PUT("/api/links/:id", h.update)
	router.DELETE("/api/links/:id", h.delete)
	router.GET("/r/:short_name", h.redirect)
}

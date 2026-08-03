package linkvisit

import (
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the metrics-related routes with the provided router.
func RegisterRoutes(shortener *shortener.Service, router gin.IRouter) {
	h := &handler{shortener: shortener}
	router.GET("/api/link_visits", h.listLinkVisits)
}

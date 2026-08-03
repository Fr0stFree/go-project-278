package linkvisit

import (
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts redirect visit routes.
func RegisterRoutes(shortener *shortener.Service, router gin.IRouter) {
	h := &handler{shortener: shortener}
	router.GET("/api/link_visits", h.listLinkVisits)
}

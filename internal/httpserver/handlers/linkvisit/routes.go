package linkvisit

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts redirect visit routes.
func RegisterRoutes(shortener shortenerService, router gin.IRouter) {
	h := &handler{shortener: shortener}
	router.GET("/api/link_visits", h.list)
}

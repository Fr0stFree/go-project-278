package linkvisit

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts redirect visit routes.
func RegisterRoutes(service Service, router gin.IRouter) {
	h := &handler{service: service}
	router.GET("/api/link_visits", h.list)
}

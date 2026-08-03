package health

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts health check routes.
func RegisterRoutes(router gin.IRouter) {
	router.GET("/ping", ping)
}

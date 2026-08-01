package health

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the health check routes with the provided router.
func RegisterRoutes(router gin.IRouter) {
	router.GET("/ping", ping)
}

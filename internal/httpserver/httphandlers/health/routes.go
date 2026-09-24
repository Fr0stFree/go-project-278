package health

import (
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts health check routes.
func RegisterRoutes(checker ReadinessChecker, healthcheckTimeout time.Duration, router gin.IRouter) {
	h := &handler{checker: checker, healthcheckTimeout: healthcheckTimeout}
	router.GET("/ping", h.ping)
	router.GET("/health", h.health)
}

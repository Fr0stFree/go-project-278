package metrics

import (
	"shortener/internal/services/metrics"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the metrics-related routes with the provided router.
func RegisterRoutes(metrics *metrics.Service, router gin.IRouter) {
	h := &handler{metrics: metrics}
	router.GET("/api/link_visits", h.listLinkVisits)
}

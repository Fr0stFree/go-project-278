package app

import (
	"net/http"
	"shortener/internal/config"

	"shortener/internal/httpserver"
	"shortener/internal/httpserver/handlers/health"
	"shortener/internal/httpserver/handlers/link"
	"shortener/internal/httpserver/handlers/metrics"
)

func prepareServer(cfg *config.HTTP, services combinedServices) *http.Server {
	router := httpserver.NewRouter()
	health.RegisterRoutes(router)
	link.RegisterRoutes(services.shortener, services.metrics, router)
	metrics.RegisterRoutes(services.metrics, router)
	server := httpserver.New(cfg, router)

	return server
}

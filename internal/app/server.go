package app

import (
	"net/http"
	"shortener/internal/config"

	"shortener/internal/httpserver"
	"shortener/internal/httpserver/handlers/health"
	"shortener/internal/httpserver/handlers/link"
)

func prepareServer(cfg *config.HTTP, services combinedServices) *http.Server {
	router := httpserver.NewRouter()
	health.RegisterRoutes(router)
	link.RegisterRoutes(services.Shortener, router)
	server := httpserver.New(cfg, router)
	return server
}
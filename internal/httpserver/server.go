// Package httpserver builds the HTTP router and server.
package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"shortener/internal/config"
	"shortener/internal/httpserver/handlers/health"
	"shortener/internal/httpserver/handlers/link"
	"shortener/internal/httpserver/handlers/linkvisit"

	"github.com/gin-gonic/gin"
)

type service interface {
	linkvisit.Service
	link.Service
}

// New creates an HTTP server for the provided handler and configuration.
func New(service service, cfg *config.HTTP) *http.Server {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health.RegisterRoutes(router)
	link.RegisterRoutes(service, router)
	linkvisit.RegisterRoutes(service, router)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	slog.Info("HTTP server configured successfully",
		slog.Int("port", cfg.Port),
		slog.Duration("read_timeout", cfg.ReadTimeout),
		slog.Duration("write_timeout", cfg.WriteTimeout),
		slog.Duration("idle_timeout", cfg.IdleTimeout),
	)

	return server
}

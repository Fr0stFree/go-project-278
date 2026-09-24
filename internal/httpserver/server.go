// Package httpserver builds the HTTP router and server.
package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"shortener/internal/config"
	"shortener/internal/httpserver/httphandlers/health"
	"shortener/internal/httpserver/httphandlers/link"
	"shortener/internal/httpserver/httphandlers/linkvisit"
	"shortener/internal/httpserver/httptools/middleware"
	"shortener/internal/integrations/sentry"
)

type service interface {
	linkvisit.Service
	link.Service
}

// New creates an HTTP server for the provided handler and configuration.
func New(service service, checker health.ReadinessChecker, cfg *config.HTTP, baseURL string) *http.Server {
	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.AccessLogger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(cfg.CORS))
	router.Use(middleware.MaxBodySize(cfg.MaxBodySize))
	router.Use(sentry.Middleware(cfg.Sentry))

	router.NoRoute(handleRouteNotFound)
	router.HandleMethodNotAllowed = true
	router.NoMethod(handleMethodNotAllowed)

	health.RegisterRoutes(checker, cfg.HealthCheckTimeout, router)
	link.RegisterRoutes(service, router, baseURL)
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

func handleRouteNotFound(ctx *gin.Context) {
	ctx.JSON(
		http.StatusNotFound,
		gin.H{"error": fmt.Sprintf("route not found: %s", ctx.Request.URL.Path)},
	)
}

func handleMethodNotAllowed(ctx *gin.Context) {
	ctx.JSON(
		http.StatusMethodNotAllowed,
		gin.H{"error": fmt.Sprintf("method %s is not allowed for %s", ctx.Request.Method, ctx.Request.URL.Path)},
	)
}

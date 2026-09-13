// Package httpserver builds the HTTP router and server.
package httpserver

import (
	"fmt"
	"net/http"
	"shortener/internal/config"
	"shortener/internal/httpserver/handlers/health"
	"shortener/internal/httpserver/handlers/link"
	"shortener/internal/httpserver/handlers/linkvisit"
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

// New creates an HTTP server for the provided handler and configuration.
func New(shortener *shortener.Service, cfg *config.HTTP) *http.Server {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health.RegisterRoutes(router)
	link.RegisterRoutes(shortener, router)
	linkvisit.RegisterRoutes(shortener, router)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return server
}

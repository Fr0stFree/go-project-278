// Package httpserver builds the HTTP router and server.
package httpserver

import (
	"fmt"
	"net/http"
	"shortener/internal/config"
)

// New creates an HTTP server for the provided handler and configuration.
func New(cfg *config.HTTP, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}

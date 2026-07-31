// Package httpserver provides functionality to create and run an HTTP server for the application.
package httpserver

import (
	"fmt"
	"net/http"
	"shortener/internal/config"
)

// New creates a new HTTP server with the specified configuration.
func New(cfg *config.HTTP, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}

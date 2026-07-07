// Package httpserver provides functionality to create and run an HTTP server for the application.
package httpserver

import (
	"fmt"
	"net/http"
	"shortener/internal/config"
)

// New creates a new HTTP server with the specified configuration.
func New(cfg config.HTTPConfig) *http.Server {
	router := newRouter()

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}

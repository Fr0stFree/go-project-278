// Package httpserver provides functionality to create and run an HTTP server for the application.
package httpserver

import (
	"fmt"
	"net/http"
	"shortener/internal/config"
	"shortener/internal/httpserver/handler"
	"shortener/internal/shortener"
)

// New creates a new HTTP server with the specified configuration.
func New(cfg config.HTTPConfig, service *shortener.Service) *http.Server {
	handler := handler.New(service)
	router := newRouter(handler)
	address := fmt.Sprintf(":%d", cfg.Port)

	return &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}

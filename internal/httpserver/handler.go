// Package handler provides HTTP handlers for the application.
package httpserver

import (
	"shortener/internal/httpserver/handlers/health"
	"shortener/internal/httpserver/handlers/link"
	"shortener/internal/shortener"
)

// Handler provides HTTP handlers for the whole application.
type Handler struct {
	Health *health.Handler
	Link   *link.Handler
}

func newHandler(shortener *shortener.Service) *Handler {
	return &Handler{
		Health: health.NewHandler(),
		Link:   link.NewHandler(shortener),
	}
}

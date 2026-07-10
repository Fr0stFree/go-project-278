// Package handler provides HTTP handlers for the application.
package handler

import (
	"shortener/internal/shortener"
)

// Handler provides HTTP handlers for the whole application.
type Handler struct {
	Health *healthHandler
	URL    *urlHandler
}

// New creates a new Handler.
func New(service *shortener.Service) *Handler {
	return &Handler{
		Health: newHealthHandler(),
		URL:    newURLHandler(service),
	}
}

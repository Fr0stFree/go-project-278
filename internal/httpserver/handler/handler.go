// Package handler provides HTTP handlers for the application.
package handler

import (
	"shortener/internal/shortener"
)

// Handler provides HTTP handlers for the whole application.
type Handler struct {
	Health *healthHandler
	Link   *linkHandler
}

// New creates a new Handler.
func New(shortener *shortener.Service) *Handler {
	return &Handler{
		Health: newHealthHandler(),
		Link:   newLinkHandler(shortener),
	}
}

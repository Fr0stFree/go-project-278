package httpserver

import (
	"shortener/internal/httpserver/handlers/health"
	"shortener/internal/httpserver/handlers/link"
	"shortener/internal/shortener"
)

type combinedHandlers struct {
	Health *health.Handler
	Link   *link.Handler
}

func newHandlers(shortener *shortener.Service) *combinedHandlers {
	return &combinedHandlers{
		Health: health.NewHandler(),
		Link:   link.NewHandler(shortener),
	}
}

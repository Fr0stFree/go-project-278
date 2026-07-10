package handler

import (
	"shortener/internal/shortener"
)

type urlHandler struct {
	Service *shortener.Service
}

func newURLHandler(service *shortener.Service) *urlHandler {
	return &urlHandler{
		Service: service,
	}
}

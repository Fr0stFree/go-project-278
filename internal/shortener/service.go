// Package shortener provides the service layer for the URL shortening application.
package shortener

import (
	"shortener/internal/storage"
)

// Service provides methods to shorten URLs and retrieve original URLs.
type Service struct {
	storage storage.Storage
}

// NewService creates a new instance of the Service with the provided storage implementation.
func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

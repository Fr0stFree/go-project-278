// Package memory provides an in-memory implementation of the Storage interface.
package memory

import (
	"shortener/internal/storage"
)

// Storage is an in-memory implementation of the Storage interface.
type Storage struct {
	data map[string]string
}

// NewStorage creates a new instance of the in-memory storage.
func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

// SaveURL saves the original URL and its corresponding short URL in memory.
func (s *Storage) SaveURL(originalURL, shortURL string) error {
	s.data[shortURL] = originalURL

	return nil
}

// GetOriginalURL retrieves the original URL corresponding to the given short URL from memory.
func (s *Storage) GetOriginalURL(shortURL string) (string, error) {
	originalURL, exists := s.data[shortURL]
	if !exists {
		return "", storage.ErrURLNotFound
	}

	return originalURL, nil
}

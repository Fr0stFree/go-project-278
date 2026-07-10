// Package storage defines the interface, errors and common types for URL storage systems.
package storage

// Storage defines the interface for a URL storage system.
type Storage interface {
	SaveURL(originalURL, shortURL string) error
	GetOriginalURL(shortURL string) (string, error)
}

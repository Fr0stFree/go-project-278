// Package shortener provides the service layer for the URL shortening application.
package shortener

import (
	"shortener/internal/storage"
)

// Service provides methods to shorten URLs and retrieve original URLs.
type Service struct {
	storage storage.Storage
	baseURL string
}

// NewService creates a new instance of the Service with the provided storage implementation.
func NewService(storage storage.Storage, baseURL string) *Service {
	return &Service{
		storage: storage,
		baseURL: baseURL,
	}
}

// ShortenLink generates a short code for the given original URL and saves the mapping in storage.
func (s *Service) ShortenLink(originalURL, shortName string) (Link, error) {
	if shortName == "" {
		shortName = toHashString(originalURL, 6)
	}

	linkDBOut, err := s.storage.SaveLink(storage.LinkDBIn{
		OriginalURL: originalURL,
		ShortName:   shortName,
	})
	if err != nil {
		return Link{}, err
	}

	return Link{
		ID:          linkDBOut.ID,
		OriginalURL: linkDBOut.OriginalURL,
		ShortName:   linkDBOut.ShortName,
		ShortURL:    s.baseURL + "/" + linkDBOut.ShortName,
	}, nil
}

// GetOriginalLink retrieves the original URL corresponding to the given short URL.
func (s *Service) GetOriginalLink(id int) (Link, error) {
	linkDBOut, err := s.storage.GetLink(id)
	if err != nil {
		return Link{}, err
	}

	return Link{
		ID:          linkDBOut.ID,
		OriginalURL: linkDBOut.OriginalURL,
		ShortName:   linkDBOut.ShortName,
		ShortURL:    s.baseURL + "/" + linkDBOut.ShortName,
	}, nil
}

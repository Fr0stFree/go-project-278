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

// CreateLink generates a short code for the given original URL and saves the mapping in storage.
func (s *Service) CreateLink(originalURL, shortName string) (Link, error) {
	if shortName == "" {
		shortName = toHashString(originalURL, 6)
	}
	linkDBIn := storage.LinkDBIn{
		OriginalURL: originalURL,
		ShortName:   shortName,
	}
	linkDBOut, err := s.storage.SaveLink(linkDBIn)
	if err != nil {
		return Link{}, err
	}

	return s.buildLink(linkDBOut), nil
}

// GetLink retrieves the original URL corresponding to the given short URL.
func (s *Service) GetLink(id int) (Link, error) {
	linkDBOut, err := s.storage.GetLink(id)
	if err != nil {
		return Link{}, err
	}

	return s.buildLink(linkDBOut), nil
}

// ListLinks retrieves a list of all shortened links stored in the service.
func (s *Service) ListLinks() ([]Link, error) {
	linksDBOut, err := s.storage.ListLinks()
	if err != nil {
		return nil, err
	}

	links := make([]Link, len(linksDBOut))
	for idx, linkDBOut := range linksDBOut {
		links[idx] = s.buildLink(linkDBOut)
	}

	return links, nil
}

func (s *Service) UpdateLink(id int, originalURL, shortName string) (Link, error) {
	linkDBIn := storage.LinkDBIn{
		OriginalURL: originalURL,
		ShortName:   shortName,
	}

	linkDBOut, err := s.storage.UpdateLink(id, linkDBIn)
	if err != nil {
		return Link{}, err
	}

	return s.buildLink(linkDBOut), nil
}

func (s *Service) DeleteLink(id int) error {
	return s.storage.DeleteLink(id)
}

func (s *Service) buildLink(linkDBOut storage.LinkDBOut) Link {
	return Link{
		ID:          linkDBOut.ID,
		OriginalURL: linkDBOut.OriginalURL,
		ShortName:   linkDBOut.ShortName,
		ShortURL:    s.baseURL + "/" + linkDBOut.ShortName,
	}
}

// Package shortener provides the service layer for the URL shortening application.
package shortener

import (
	"shortener/internal/config"
	"shortener/internal/storage"
)

// Service provides methods to shorten URLs and retrieve original URLs.
type Service struct {
	linkRepository storage.AbstractLinkRepository
	config         *config.App
}

// NewService creates a new instance of the Service with the provided storage implementation.
func NewService(linkRepository storage.AbstractLinkRepository, config *config.App) *Service {
	return &Service{
		linkRepository: linkRepository,
		config:         config,
	}
}

// CreateLink generates a short code for the given original URL and saves the mapping in storage.
func (s *Service) CreateLink(originalURL, shortName string) (Link, error) {
	if shortName == "" {
		shortName = toHashString(originalURL, 6)
	}

	insert := storage.LinkInsert{
		OriginalURL: originalURL,
		ShortName:   shortName,
	}

	record, err := s.linkRepository.SaveLink(insert)
	if err != nil {
		return Link{}, err
	}

	return s.buildLink(record), nil
}

// GetLink retrieves the original URL corresponding to the given short URL.
func (s *Service) GetLink(id int) (Link, error) {
	record, err := s.linkRepository.GetLink(id)
	if err != nil {
		return Link{}, err
	}

	return s.buildLink(record), nil
}

// ListLinks retrieves a list of all shortened links stored in the service.
func (s *Service) ListLinks() ([]Link, error) {
	records, err := s.linkRepository.ListLinks()
	if err != nil {
		return nil, err
	}

	links := make([]Link, len(records))
	for idx, record := range records {
		links[idx] = s.buildLink(record)
	}

	return links, nil
}

func (s *Service) UpdateLink(id int, originalURL, shortName string) (Link, error) {
	update := storage.LinkUpdate{
		OriginalURL: originalURL,
		ShortName:   shortName,
	}

	record, err := s.linkRepository.UpdateLink(id, update)
	if err != nil {
		return Link{}, err
	}

	return s.buildLink(record), nil
}

func (s *Service) DeleteLink(id int) error {
	return s.linkRepository.DeleteLink(id)
}

func (s *Service) buildLink(record storage.LinkRecord) Link {
	return Link{
		ID:          record.ID,
		OriginalURL: record.OriginalURL,
		ShortName:   record.ShortName,
		ShortURL:    s.config.BaseURL + "/" + record.ShortName,
	}
}

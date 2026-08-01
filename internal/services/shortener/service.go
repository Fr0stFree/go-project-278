// Package shortener provides the service layer for the URL shortening application.
package shortener

import (
	"shortener/internal/config"
	"shortener/internal/database/repositories/link"
)

// Service provides methods to shorten URLs and retrieve original URLs.
type Service struct {
	repo   link.AbstractRepository
	config *config.App
}

// NewService creates a new instance of the Service with the provided storage implementation.
func NewService(linkRepository link.AbstractRepository, config *config.App) *Service {
	return &Service{
		repo:   linkRepository,
		config: config,
	}
}

// CreateLink generates a short code for the given original URL and saves the mapping in storage.
func (s *Service) CreateLink(originalURL, shortName string) (Link, error) {
	if shortName == "" {
		shortName = toHashString(originalURL, 6)
	}

	insert := link.Insert{
		OriginalURL: originalURL,
		ShortName:   shortName,
	}

	record, err := s.repo.CreateOne(insert)
	if err != nil {
		return Link{}, err
	}

	return s.buildLink(record), nil
}

// GetLink retrieves the original URL corresponding to the given short URL.
func (s *Service) GetLink(id int) (Link, error) {
	record, err := s.repo.GetByID(id)
	if err != nil {
		return Link{}, err
	}

	return s.buildLink(record), nil
}

// ListLinks retrieves a list of all shortened links stored in the service.
func (s *Service) ListLinks(linkRange *LinkRange) ([]Link, error) {
	options := link.NewFilterOpts()

	if linkRange != nil {
		options.WithRange(linkRange.From, linkRange.To)
	}

	records, err := s.repo.GetMany(options)
	if err != nil {
		return nil, err
	}

	links := make([]Link, len(records))
	for idx, record := range records {
		links[idx] = s.buildLink(record)
	}

	return links, nil
}

// UpdateLink updates the original URL and/or short name for the link with the specified ID.
func (s *Service) UpdateLink(id int, originalURL, shortName string) (Link, error) {
	update := link.Update{
		OriginalURL: originalURL,
		ShortName:   shortName,
	}

	record, err := s.repo.UpdateByID(id, update)
	if err != nil {
		return Link{}, err
	}

	return s.buildLink(record), nil
}

// DeleteLink removes the link with the specified ID from the storage.
func (s *Service) DeleteLink(id int) error {
	return s.repo.DeleteByID(id)
}

func (s *Service) buildLink(record link.Record) Link {
	return Link{
		ID:          record.ID,
		OriginalURL: record.OriginalURL,
		ShortName:   record.ShortName,
		ShortURL:    s.config.BaseURL + "/" + record.ShortName,
	}
}

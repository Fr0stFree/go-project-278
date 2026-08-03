// Package shortener provides the service layer for the URL shortening application.
package shortener

import (
	"shortener/internal/config"
	"shortener/internal/database/repositories/link"
	"shortener/internal/database/repositories/linkvisit"
	"time"
)

// Service provides methods to shorten URLs and retrieve original URLs.
type Service struct {
	linkRepo      link.AbstractRepository
	linkVisitRepo linkvisit.AbstractRepository
	cfg           *config.App
}

// NewService creates a new instance of the Service with the provided storage implementation.
func NewService(linkRepository link.AbstractRepository, linkVisitRepository linkvisit.AbstractRepository, config *config.App) *Service {
	return &Service{
		linkRepo:      linkRepository,
		linkVisitRepo: linkVisitRepository,
		cfg:           config,
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

	record, err := s.linkRepo.CreateOne(insert)
	if err != nil {
		return Link{}, err
	}

	return s.buildLink(record), nil
}

// GetLink retrieves the original URL corresponding to the given short URL.
func (s *Service) GetLink(id uint) (Link, error) {
	record, err := s.linkRepo.GetByID(id)
	if err != nil {
		return Link{}, err
	}

	return s.buildLink(record), nil
}

// GetRedirectLink retrieves the original URL corresponding to the given short name for redirection purposes.
func (s *Service) GetRedirectLink(shortName string) (Link, error) {
	opts, err := link.NewFilterOpts().WithShortNames(shortName).WithRange(0, 0)
	if err != nil {
		return Link{}, err
	}

	records, err := s.linkRepo.GetMany(*opts)
	if err != nil {
		return Link{}, err
	}

	if len(records) == 0 {
		return Link{}, link.ErrNotFound
	}

	return s.buildLink(records[0]), nil
}

func (s *Service) ListLinksWithCount(opts *link.FilterOpts) ([]Link, int, error) {
	if opts == nil {
		opts = link.NewFilterOpts()
	}

	records, err := s.linkRepo.GetMany(*opts)
	if err != nil {
		return nil, 0, err
	}

	links := make([]Link, len(records))
	for idx, record := range records {
		links[idx] = s.buildLink(record)
	}

	count, err := s.linkRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return links, count, nil
}

// UpdateLink updates the original URL and/or short name for the link with the specified ID.
func (s *Service) UpdateLink(id uint, originalURL, shortName string) (Link, error) {
	update := link.Update{
		OriginalURL: originalURL,
		ShortName:   shortName,
	}

	record, err := s.linkRepo.UpdateByID(id, update)
	if err != nil {
		return Link{}, err
	}

	return s.buildLink(record), nil
}

// DeleteLink removes the link with the specified ID from the storage.
func (s *Service) DeleteLink(id uint) error {
	return s.linkRepo.DeleteByID(id)
}

func (s *Service) buildLink(record link.Record) Link {
	return Link{
		ID:          record.ID,
		OriginalURL: record.OriginalURL,
		ShortName:   record.ShortName,
		ShortURL:    s.cfg.BaseURL + "/r/" + record.ShortName,
	}
}

// SaveLinkVisit saves a new link visit record in the repository with the provided details.
func (s *Service) SaveLinkVisit(linkID uint, ip, userAgent, referrer string, status uint) (LinkVisit, error) {
	insert := linkvisit.Insert{
		LinkID:    linkID,
		IP:        ip,
		UserAgent: userAgent,
		Referrer:  referrer,
		Status:    status,
	}

	record, err := s.linkVisitRepo.CreateOne(insert)
	if err != nil {
		return LinkVisit{}, err
	}

	return s.buildLinkVisit(record), nil
}

func (s *Service) ListLinkVisitsWithCount(options *linkvisit.FilterOpts) ([]LinkVisit, int, error) {
	if options == nil {
		options = linkvisit.NewFilterOpts()
	}

	records, err := s.linkVisitRepo.GetMany(*options)
	if err != nil {
		return nil, 0, err
	}

	visits := make([]LinkVisit, len(records))
	for i, record := range records {
		visits[i] = s.buildLinkVisit(record)
	}

	count, err := s.linkVisitRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return visits, count, nil
}

func (s *Service) buildLinkVisit(record linkvisit.Record) LinkVisit {
	return LinkVisit{
		ID:        record.ID,
		LinkID:    record.LinkID,
		CreatedAt: record.CreatedAt.Format(time.RFC3339),
		UpdatedAt: record.UpdatedAt.Format(time.RFC3339),
		IP:        record.IP,
		UserAgent: record.UserAgent,
		Status:    record.Status,
		Referrer:  record.Referrer,
	}
}

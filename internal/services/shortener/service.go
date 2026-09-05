// Package shortener implements link shortening use cases.
package shortener

import (
	"errors"
	"shortener/internal/config"
	"shortener/internal/db/models"
	"shortener/internal/db/models/link"
	"shortener/internal/db/models/linkvisit"
	"time"
)

// Service coordinates link and visit repositories.
type Service struct {
	linkRepo      LinkRepository
	linkVisitRepo LinkVisitRepository
	cfg           *config.App
}

// NewService creates a shortener service with link and visit repositories.
func NewService(linkRepository LinkRepository, linkVisitRepository LinkVisitRepository, config *config.App) *Service {
	return &Service{
		linkRepo:      linkRepository,
		linkVisitRepo: linkVisitRepository,
		cfg:           config,
	}
}

// CreateLink creates a shortened link, generating a short name when one is not provided.
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
		return Link{}, s.mapStorageErrorToServiceError(err)
	}

	return s.buildLink(record), nil
}

// GetLink returns a shortened link by ID.
func (s *Service) GetLink(id uint) (Link, error) {
	record, err := s.linkRepo.GetByID(id)
	if err != nil {
		return Link{}, s.mapStorageErrorToServiceError(err)
	}

	return s.buildLink(record), nil
}

// GetRedirectLink returns a shortened link by short name.
func (s *Service) GetRedirectLink(shortName string) (Link, error) {
	builder := NewLinkListOptionsBuilder()
	builder.WithShortNames(shortName)
	builder.WithRange(0, 0)

	if builder.err != nil {
		return Link{}, builder.err
	}

	records, err := s.linkRepo.GetMany(builder.build())
	if err != nil {
		return Link{}, s.mapStorageErrorToServiceError(err)
	}

	if len(records) == 0 {
		return Link{}, NewNotFoundError("link not found")
	}

	return s.buildLink(records[0]), nil
}

// ListLinksWithCount returns filtered links and the total link count.
func (s *Service) ListLinksWithCount(builder *LinkListOptionsBuilder) ([]Link, int, error) {
	if builder == nil {
		builder = NewLinkListOptionsBuilder()
	}

	if builder.err != nil {
		return nil, 0, builder.err
	}

	records, err := s.linkRepo.GetMany(builder.build())
	if err != nil {
		return nil, 0, s.mapStorageErrorToServiceError(err)
	}

	links := make([]Link, len(records))
	for idx, record := range records {
		links[idx] = s.buildLink(record)
	}

	count, err := s.linkRepo.Count()
	if err != nil {
		return nil, 0, s.mapStorageErrorToServiceError(err)
	}

	return links, count, nil
}

// UpdateLink replaces URL fields for a shortened link by ID.
func (s *Service) UpdateLink(id uint, originalURL, shortName string) (Link, error) {
	update := link.Update{
		OriginalURL: originalURL,
		ShortName:   shortName,
	}

	record, err := s.linkRepo.UpdateByID(id, update)
	if err != nil {
		return Link{}, s.mapStorageErrorToServiceError(err)
	}

	return s.buildLink(record), nil
}

// DeleteLink removes a shortened link by ID.
func (s *Service) DeleteLink(id uint) error {
	err := s.linkRepo.DeleteByID(id)
	if err != nil {
		return s.mapStorageErrorToServiceError(err)
	}

	return nil
}

func (s *Service) buildLink(record link.Record) Link {
	return Link{
		ID:          record.ID,
		OriginalURL: record.OriginalURL,
		ShortName:   record.ShortName,
		ShortURL:    s.cfg.BaseURL + "/r/" + record.ShortName,
	}
}

// SaveLinkVisit records a redirect attempt for a shortened link.
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
		return LinkVisit{}, s.mapStorageErrorToServiceError(err)
	}

	return s.buildLinkVisit(record), nil
}

// ListLinkVisitsWithCount returns filtered visits and the total visit count.
func (s *Service) ListLinkVisitsWithCount(builder *LinkVisitListOptionsBuilder) ([]LinkVisit, int, error) {
	if builder == nil {
		builder = NewLinkVisitListOptionsBuilder()
	}

	if builder.err != nil {
		return nil, 0, builder.err
	}

	records, err := s.linkVisitRepo.GetMany(builder.build())
	if err != nil {
		return nil, 0, s.mapStorageErrorToServiceError(err)
	}

	visits := make([]LinkVisit, len(records))
	for i, record := range records {
		visits[i] = s.buildLinkVisit(record)
	}

	count, err := s.linkVisitRepo.Count()
	if err != nil {
		return nil, 0, s.mapStorageErrorToServiceError(err)
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

func (s *Service) mapStorageErrorToServiceError(err error) error {
	switch {
	case errors.Is(err, models.ErrObjectDoesNotExist):
		return NewNotFoundError("link not found")
	case errors.Is(err, models.ErrObjectAlreadyExists):
		return NewConflictError("shortname already in use", "short_name")
	default:
		return err
	}
}

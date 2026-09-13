// Package shortener implements link shortening use cases.
package shortener

import (
	"context"
	"errors"
	"shortener/internal/common/textutils"
	"shortener/internal/config"
	"shortener/internal/db/models"
	"shortener/internal/db/models/link"
	"shortener/internal/db/models/linkvisit"
	"time"

	"github.com/samber/lo"
)

type linkVisitRepository interface {
	CreateOne(ctx context.Context, insert linkvisit.Insert) (linkvisit.Record, error)
	GetMany(ctx context.Context, options linkvisit.ListOptions) ([]linkvisit.Record, error)
	Count(ctx context.Context) (int, error)
}

type linkRepository interface {
	CreateOne(ctx context.Context, insert link.Insert) (link.Record, error)
	GetByID(ctx context.Context, ID uint) (link.Record, error)
	GetMany(ctx context.Context, options link.ListOptions) ([]link.Record, error)
	Count(ctx context.Context) (int, error)
	UpdateByID(ctx context.Context, ID uint, update link.Update) (link.Record, error)
	DeleteByID(ctx context.Context, ID uint) error
}

// Service coordinates link and visit repositories.
type Service struct {
	links      linkRepository
	linkVisits linkVisitRepository
	cfg        *config.App
}

// NewService creates a shortener service with link and visit repositories.
func NewService(linkRepository linkRepository, linkVisitRepository linkVisitRepository, config *config.App) *Service {
	return &Service{
		links:      linkRepository,
		linkVisits: linkVisitRepository,
		cfg:        config,
	}
}

// CreateLink creates a shortened link, generating a short name when one is not provided.
func (s *Service) CreateLink(ctx context.Context, originalURL, shortName string) (Link, error) {
	isShortNameProvided := shortName != ""
	shortName = lo.Ternary(isShortNameProvided, shortName, textutils.RandomString(6))

	insert := link.Insert{
		OriginalURL: originalURL,
		ShortName:   shortName,
	}

	record, err := s.links.CreateOne(ctx, insert)
	if err != nil {
		if errors.Is(err, models.ErrObjectAlreadyExists) && !isShortNameProvided {
			// If the short name was generated and already exists, try again with a new random short name.
			return s.CreateLink(ctx, originalURL, "")
		}

		return Link{}, s.mapStorageErrorToServiceError(err)
	}

	return s.buildLink(record), nil
}

// GetLink returns a shortened link by ID.
func (s *Service) GetLink(ctx context.Context, id uint) (Link, error) {
	record, err := s.links.GetByID(ctx, id)
	if err != nil {
		return Link{}, s.mapStorageErrorToServiceError(err)
	}

	return s.buildLink(record), nil
}

// GetRedirectLink returns a shortened link by short name.
func (s *Service) GetRedirectLink(ctx context.Context, shortName string) (Link, error) {
	builder := NewLinkListOptionsBuilder()
	builder.WithShortNames(shortName)
	builder.WithRange(0, 0)

	if builder.err != nil {
		return Link{}, builder.err
	}

	records, err := s.links.GetMany(ctx, builder.build())
	if err != nil {
		return Link{}, s.mapStorageErrorToServiceError(err)
	}

	if len(records) == 0 {
		return Link{}, NewNotFoundError("link not found")
	}

	return s.buildLink(records[0]), nil
}

// ListLinksWithCount returns filtered links and the total link count.
func (s *Service) ListLinksWithCount(ctx context.Context, builder *LinkListOptionsBuilder) ([]Link, int, error) {
	if builder == nil {
		builder = NewLinkListOptionsBuilder()
	}

	if builder.err != nil {
		return nil, 0, builder.err
	}

	records, err := s.links.GetMany(ctx, builder.build())
	if err != nil {
		return nil, 0, s.mapStorageErrorToServiceError(err)
	}

	links := make([]Link, len(records))
	for idx, record := range records {
		links[idx] = s.buildLink(record)
	}

	count, err := s.links.Count(ctx)
	if err != nil {
		return nil, 0, s.mapStorageErrorToServiceError(err)
	}

	return links, count, nil
}

// UpdateLink replaces URL fields for a shortened link by ID.
func (s *Service) UpdateLink(ctx context.Context, id uint, originalURL, shortName string) (Link, error) {
	update := link.Update{
		OriginalURL: originalURL,
		ShortName:   shortName,
	}

	record, err := s.links.UpdateByID(ctx, id, update)
	if err != nil {
		return Link{}, s.mapStorageErrorToServiceError(err)
	}

	return s.buildLink(record), nil
}

// DeleteLink removes a shortened link by ID.
func (s *Service) DeleteLink(ctx context.Context, id uint) error {
	err := s.links.DeleteByID(ctx, id)
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
func (s *Service) SaveLinkVisit(ctx context.Context, linkID uint, ip, userAgent, referrer string, status uint) (LinkVisit, error) {
	insert := linkvisit.Insert{
		LinkID:    linkID,
		IP:        ip,
		UserAgent: userAgent,
		Referrer:  referrer,
		Status:    status,
	}

	record, err := s.linkVisits.CreateOne(ctx, insert)
	if err != nil {
		return LinkVisit{}, s.mapStorageErrorToServiceError(err)
	}

	return s.buildLinkVisit(record), nil
}

// ListLinkVisitsWithCount returns filtered visits and the total visit count.
func (s *Service) ListLinkVisitsWithCount(ctx context.Context, builder *LinkVisitListOptionsBuilder) ([]LinkVisit, int, error) {
	if builder == nil {
		builder = NewLinkVisitListOptionsBuilder()
	}

	if builder.err != nil {
		return nil, 0, builder.err
	}

	records, err := s.linkVisits.GetMany(ctx, builder.build())
	if err != nil {
		return nil, 0, s.mapStorageErrorToServiceError(err)
	}

	visits := make([]LinkVisit, len(records))
	for i, record := range records {
		visits[i] = s.buildLinkVisit(record)
	}

	count, err := s.linkVisits.Count(ctx)
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

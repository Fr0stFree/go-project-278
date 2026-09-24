// Package shortener implements link shortening use cases.
package shortener

import (
	"context"
	"errors"
	"regexp"

	"shortener/internal/common/textutils"
)

type linkVisitRepository interface {
	CreateOne(ctx context.Context, params CreateLinkVisitParams) (LinkVisit, error)
	GetMany(ctx context.Context, options LinkVisitListOptions) ([]LinkVisit, error)
	Count(ctx context.Context) (int, error)
}

type linkRepository interface {
	CreateOne(ctx context.Context, params CreateLinkParams) (Link, error)
	GetByID(ctx context.Context, ID uint) (Link, error)
	GetMany(ctx context.Context, options LinkListOptions) ([]Link, error)
	Count(ctx context.Context) (int, error)
	UpdateByID(ctx context.Context, ID uint, params UpdateLinkParams) (Link, error)
	DeleteByID(ctx context.Context, ID uint) error
}

type serviceOpts struct {
	shortNameGenerationMaxAttempts int
	shortNameDefaultLength         int
	shortNameMinLength             int
	shortNameMaxLength             int
	shortNameRegexPattern          regexp.Regexp
}

// Service coordinates link shortening and visit tracking operations.
type Service struct {
	links      linkRepository
	linkVisits linkVisitRepository
	opts       serviceOpts
}

// NewService creates a shortener service with link and visit repositories.
func NewService(
	linkRepository linkRepository,
	linkVisitRepository linkVisitRepository,
) *Service {
	return &Service{
		links:      linkRepository,
		linkVisits: linkVisitRepository,
		opts: serviceOpts{
			shortNameGenerationMaxAttempts: 10,
			shortNameDefaultLength:         6,
			shortNameMinLength:             3,
			shortNameMaxLength:             32,
			shortNameRegexPattern:          *regexp.MustCompile(`^[a-zA-Z0-9_-]+$`),
		},
	}
}

// CreateLink creates a shortened link, generating a short name when one is not provided.
func (s *Service) CreateLink(ctx context.Context, originalURL, shortName string) (Link, error) {
	isShortNameProvided := shortName != ""

	for range s.opts.shortNameGenerationMaxAttempts {
		if !isShortNameProvided {
			shortName = textutils.RandomString(s.opts.shortNameDefaultLength)
		}

		if err := s.validateShortName(shortName); err != nil {
			return Link{}, NewValidationError(err.Error(), "short_name")
		}

		params := CreateLinkParams{
			OriginalURL: originalURL,
			ShortName:   shortName,
		}

		link, err := s.links.CreateOne(ctx, params)
		if err != nil {
			if errors.Is(err, ErrRepositoryConflict) && !isShortNameProvided {
				continue
			}

			return Link{}, s.mapRepositoryError(err)
		}

		return link, nil
	}

	return Link{}, ErrShortNameGenerationExhausted
}

// GetLink returns a shortened link by ID.
func (s *Service) GetLink(ctx context.Context, id uint) (Link, error) {
	link, err := s.links.GetByID(ctx, id)
	if err != nil {
		return Link{}, s.mapRepositoryError(err)
	}

	return link, nil
}

// GetRedirectLink returns a shortened link by short name.
func (s *Service) GetRedirectLink(ctx context.Context, shortName string) (Link, error) {
	builder := NewLinkListOptionsBuilder()
	builder.WithShortNames(shortName)
	builder.WithRange(0, 0)

	if builder.err != nil {
		return Link{}, builder.err
	}

	links, err := s.links.GetMany(ctx, builder.build())
	if err != nil {
		return Link{}, s.mapRepositoryError(err)
	}

	if len(links) == 0 {
		return Link{}, NewNotFoundError("link not found")
	}

	return links[0], nil
}

// ListLinksWithCount returns filtered links and the total link count.
func (s *Service) ListLinksWithCount(ctx context.Context, builder *LinkListOptionsBuilder) ([]Link, int, error) {
	if builder == nil {
		builder = NewLinkListOptionsBuilder()
	}

	if builder.err != nil {
		return nil, 0, builder.err
	}

	links, err := s.links.GetMany(ctx, builder.build())
	if err != nil {
		return nil, 0, s.mapRepositoryError(err)
	}

	count, err := s.links.Count(ctx)
	if err != nil {
		return nil, 0, s.mapRepositoryError(err)
	}

	return links, count, nil
}

// UpdateLink replaces URL fields for a shortened link by ID.
func (s *Service) UpdateLink(ctx context.Context, id uint, originalURL, shortName string) (Link, error) {
	if err := s.validateShortName(shortName); err != nil {
		return Link{}, NewValidationError(err.Error(), "short_name")
	}

	params := UpdateLinkParams{
		OriginalURL: originalURL,
		ShortName:   shortName,
	}

	link, err := s.links.UpdateByID(ctx, id, params)
	if err != nil {
		return Link{}, s.mapRepositoryError(err)
	}

	return link, nil
}

// DeleteLink removes a shortened link by ID.
func (s *Service) DeleteLink(ctx context.Context, id uint) error {
	err := s.links.DeleteByID(ctx, id)
	if err != nil {
		return s.mapRepositoryError(err)
	}

	return nil
}

// SaveLinkVisit records a redirect attempt for a shortened link.
func (s *Service) SaveLinkVisit(ctx context.Context, linkID uint, ip, userAgent, referrer string, status uint) (LinkVisit, error) {
	params := CreateLinkVisitParams{
		LinkID:    linkID,
		IP:        ip,
		UserAgent: userAgent,
		Referrer:  referrer,
		Status:    status,
	}

	visit, err := s.linkVisits.CreateOne(ctx, params)
	if err != nil {
		return LinkVisit{}, s.mapRepositoryError(err)
	}

	return visit, nil
}

// ListLinkVisitsWithCount returns filtered visits and the total visit count.
func (s *Service) ListLinkVisitsWithCount(ctx context.Context, builder *LinkVisitListOptionsBuilder) ([]LinkVisit, int, error) {
	if builder == nil {
		builder = NewLinkVisitListOptionsBuilder()
	}

	if builder.err != nil {
		return nil, 0, builder.err
	}

	visits, err := s.linkVisits.GetMany(ctx, builder.build())
	if err != nil {
		return nil, 0, s.mapRepositoryError(err)
	}

	count, err := s.linkVisits.Count(ctx)
	if err != nil {
		return nil, 0, s.mapRepositoryError(err)
	}

	return visits, count, nil
}

func (s *Service) mapRepositoryError(err error) error {
	switch {
	case errors.Is(err, ErrRepositoryNotFound):
		return NewNotFoundError("link not found")
	case errors.Is(err, ErrRepositoryConflict):
		return NewConflictError("shortname already in use", "short_name")
	default:
		return err
	}
}

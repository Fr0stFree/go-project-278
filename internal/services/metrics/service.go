// Package metrics provides services for managing link visits and interacting with the link visit repository.
package metrics

import (
	"shortener/internal/config"
	"shortener/internal/database/repositories/linkvisit"
	"time"
)

// Service provides methods for managing link visits and interacting with the link visit repository.
type Service struct {
	repo linkvisit.AbstractRepository
	cfg  *config.App
}

// NewService creates a new instance of the Service with the provided link visit repository and configuration.
func NewService(repo linkvisit.AbstractRepository, cfg *config.App) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

// SaveLinkVisit saves a new link visit record in the repository with the provided details.
func (s *Service) SaveLinkVisit(linkID int, ip, userAgent, referrer string, status int) (LinkVisit, error) {
	insert := linkvisit.Insert{
		LinkID:    linkID,
		IP:        ip,
		UserAgent: userAgent,
		Referrer:  referrer,
		Status:    status,
	}

	record, err := s.repo.CreateOne(insert)
	if err != nil {
		return LinkVisit{}, err
	}

	return s.buildLinkVisit(record), nil
}

// ListLinkVisits retrieves a list of link visit records from the repository based on the provided filter options.
func (s *Service) ListLinkVisits(options *linkvisit.FilterOpts) ([]LinkVisit, error) {
	if options == nil {
		options = linkvisit.NewFilterOpts()
	}

	records, err := s.repo.GetMany(*options)
	if err != nil {
		return nil, err
	}

	visits := make([]LinkVisit, len(records))
	for i, record := range records {
		visits[i] = s.buildLinkVisit(record)
	}

	return visits, nil
}

// CountLinkVisits returns the total number of link visit records in the repository.
func (s *Service) CountLinkVisits() (int, error) {
	count, err := s.repo.Count()
	if err != nil {
		return 0, err
	}

	return count, nil
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

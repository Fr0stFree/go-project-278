// Package linkvisit stores redirect visit records.
package linkvisit

import (
	"context"
	"fmt"
	"shortener/internal/services/shortener"

	"gorm.io/gorm"
)

// Repository stores redirect visits in PostgreSQL.
type Repository struct {
	*gorm.DB
}

// NewRepository creates a visit repository backed by the provided database.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

// CreateOne inserts a redirect visit row.
func (r *Repository) CreateOne(ctx context.Context, params shortener.CreateLinkVisitParams) (shortener.LinkVisit, error) {
	record := Record{
		LinkID:    params.LinkID,
		IP:        params.IP,
		UserAgent: params.UserAgent,
		Status:    params.Status,
		Referrer:  params.Referrer,
	}

	result := r.DB.WithContext(ctx).Create(&record)
	if result.Error != nil {
		return shortener.LinkVisit{}, result.Error
	}

	return toServiceLinkVisit(record), nil
}

// GetMany returns visit rows matching the provided list options.
func (r *Repository) GetMany(ctx context.Context, options shortener.LinkVisitListOptions) ([]shortener.LinkVisit, error) {
	records := make([]Record, 0)
	statement := r.DB.WithContext(ctx).Model(&Record{})

	if len(options.LinkIDs) > 0 {
		statement = statement.Where("link_id IN ?", options.LinkIDs)
	}

	result := statement.
		Limit(options.Limit).
		Offset(options.Offset).
		Order(fmt.Sprintf("%s %s", linkVisitSortColumn(options.SortBy), options.SortOrder)).
		Find(&records)

	if result.Error != nil {
		return nil, result.Error
	}

	visits := make([]shortener.LinkVisit, len(records))
	for i, record := range records {
		visits[i] = toServiceLinkVisit(record)
	}

	return visits, nil
}

func toServiceLinkVisit(record Record) shortener.LinkVisit {
	return shortener.LinkVisit{
		ID: record.ID, LinkID: record.LinkID, CreatedAt: record.CreatedAt, UpdatedAt: record.UpdatedAt,
		IP: record.IP, UserAgent: record.UserAgent, Status: record.Status, Referrer: record.Referrer,
	}
}

func linkVisitSortColumn(field shortener.LinkVisitSortField) string {
	switch field {
	case shortener.LinkVisitSortByLinkID:
		return "link_id"
	case shortener.LinkVisitSortByCreatedAt:
		return "created_at"
	default:
		return "id"
	}
}

// Count returns the total number of visit rows.
func (r *Repository) Count(ctx context.Context) (int, error) {
	var count int64

	result := r.DB.WithContext(ctx).Model(&Record{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

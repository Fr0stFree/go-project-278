package link

import (
	"context"
	"errors"
	"fmt"
	"shortener/internal/services/shortener"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository stores shortened links in PostgreSQL.
type Repository struct {
	*gorm.DB
}

// NewRepository creates a link repository backed by the provided database.
func NewRepository(database *gorm.DB) *Repository {
	return &Repository{database}
}

// CreateOne inserts a shortened link row.
func (r *Repository) CreateOne(ctx context.Context, params shortener.CreateLinkParams) (shortener.Link, error) {
	record := Record{
		OriginalURL: params.OriginalURL,
		ShortName:   params.ShortName,
	}

	result := r.DB.WithContext(ctx).Create(&record)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return shortener.Link{}, shortener.ErrRepositoryConflict
		}

		return shortener.Link{}, result.Error
	}

	return toServiceLink(record), nil
}

// GetByID returns a link row by ID.
func (r *Repository) GetByID(ctx context.Context, ID uint) (shortener.Link, error) {
	var record Record

	result := r.DB.WithContext(ctx).First(&record, ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return shortener.Link{}, shortener.ErrRepositoryNotFound
		}

		return shortener.Link{}, result.Error
	}

	return toServiceLink(record), nil
}

// GetMany returns link rows matching the provided list options.
func (r *Repository) GetMany(ctx context.Context, options shortener.LinkListOptions) ([]shortener.Link, error) {
	records := make([]Record, 0)

	statement := r.DB.WithContext(ctx).Model(&Record{})
	if len(options.ShortNames) > 0 {
		statement = statement.Where("short_name IN ?", options.ShortNames)
	}

	result := statement.
		Limit(options.Limit).
		Offset(options.Offset).
		Order(fmt.Sprintf("%s %s", linkSortColumn(options.SortBy), options.SortOrder)).
		Find(&records)

	if result.Error != nil {
		return nil, result.Error
	}

	links := make([]shortener.Link, len(records))
	for i, record := range records {
		links[i] = toServiceLink(record)
	}

	return links, nil
}

// Count returns the total number of link rows.
func (r *Repository) Count(ctx context.Context) (int, error) {
	var count int64

	result := r.DB.WithContext(ctx).Model(&Record{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

// UpdateByID replaces URL fields for a link row by ID.
func (r *Repository) UpdateByID(ctx context.Context, ID uint, params shortener.UpdateLinkParams) (shortener.Link, error) {
	var record Record

	result := r.DB.
		WithContext(ctx).
		Model(&record).
		Clauses(clause.Returning{}).
		Where("id = ?", ID).
		Updates(Record{OriginalURL: params.OriginalURL, ShortName: params.ShortName})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return shortener.Link{}, shortener.ErrRepositoryConflict
		}

		return shortener.Link{}, result.Error
	}

	if result.RowsAffected == 0 {
		return shortener.Link{}, shortener.ErrRepositoryNotFound
	}

	return toServiceLink(record), nil
}

// DeleteByID deletes a link row by ID.
func (r *Repository) DeleteByID(ctx context.Context, ID uint) error {
	result := r.DB.WithContext(ctx).Where("id = ?", ID).Delete(&Record{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return shortener.ErrRepositoryNotFound
	}

	return nil
}

func toServiceLink(record Record) shortener.Link {
	return shortener.Link{
		ID: record.ID, OriginalURL: record.OriginalURL, ShortName: record.ShortName,
		CreatedAt: record.CreatedAt, UpdatedAt: record.UpdatedAt,
	}
}

func linkSortColumn(field shortener.LinkSortField) string {
	switch field {
	case shortener.LinkSortByOriginalURL:
		return "original_url"
	case shortener.LinkSortByShortName:
		return "short_name"
	case shortener.LinkSortByCreatedAt:
		return "created_at"
	default:
		return "id"
	}
}

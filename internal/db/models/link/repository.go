package link

import (
	"context"
	"errors"
	"fmt"
	"shortener/internal/db/storage"

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
func (r *Repository) CreateOne(ctx context.Context, insert Insert) (Record, error) {
	record := Record{
		OriginalURL: insert.OriginalURL,
		ShortName:   insert.ShortName,
	}

	result := r.DB.WithContext(ctx).Create(&record)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return Record{}, storage.ErrObjectAlreadyExists
		}

		return Record{}, result.Error
	}

	return record, nil
}

// GetByID returns a link row by ID.
func (r *Repository) GetByID(ctx context.Context, ID uint) (Record, error) {
	var record Record

	result := r.DB.WithContext(ctx).First(&record, ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Record{}, storage.ErrObjectDoesNotExist
		}

		return Record{}, result.Error
	}

	return record, nil
}

// GetMany returns link rows matching the provided list options.
func (r *Repository) GetMany(ctx context.Context, options ListOptions) ([]Record, error) {
	records := make([]Record, 0)

	statement := r.DB.WithContext(ctx).Model(&Record{})
	if len(options.ShortNames) > 0 {
		statement = statement.Where("short_name IN ?", options.ShortNames)
	}

	result := statement.
		Limit(options.Limit).
		Offset(options.Offset).
		Order(fmt.Sprintf("%s %s", options.SortBy, options.SortOrder)).
		Find(&records)

	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
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
func (r *Repository) UpdateByID(ctx context.Context, ID uint, update Update) (Record, error) {
	var record Record

	result := r.DB.
		WithContext(ctx).
		Model(&record).
		Clauses(clause.Returning{}).
		Where("id = ?", ID).
		Updates(Record{OriginalURL: update.OriginalURL, ShortName: update.ShortName})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return Record{}, storage.ErrObjectAlreadyExists
		}

		return Record{}, result.Error
	}

	if result.RowsAffected == 0 {
		return Record{}, storage.ErrObjectDoesNotExist
	}

	return record, nil
}

// DeleteByID deletes a link row by ID.
func (r *Repository) DeleteByID(ctx context.Context, ID uint) error {
	result := r.DB.WithContext(ctx).Where("id = ?", ID).Delete(&Record{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return storage.ErrObjectDoesNotExist
	}

	return nil
}

// Package linkvisit stores redirect visit records.
package linkvisit

import (
	"fmt"
	"shortener/internal/db"
)

// Repository stores redirect visits in PostgreSQL.
type Repository struct {
	*db.DataBase
}

// NewRepository creates a visit repository backed by the provided database.
func NewRepository(db *db.DataBase) *Repository {
	return &Repository{db}
}

// CreateOne inserts a redirect visit row.
func (r *Repository) CreateOne(insert Insert) (Record, error) {
	record := Record{
		LinkID:    insert.LinkID,
		IP:        insert.IP,
		UserAgent: insert.UserAgent,
		Status:    insert.Status,
		Referrer:  insert.Referrer,
	}

	result := r.DB.Create(&record)
	if result.Error != nil {
		return Record{}, result.Error
	}

	return record, nil
}

// GetMany returns visit rows matching the provided list options.
func (r *Repository) GetMany(options ListOptions) ([]Record, error) {
	records := make([]Record, 0)

	result := r.DB.
		Limit(options.Limit).
		Offset(options.Offset).
		Order(fmt.Sprintf("%s %s", options.SortBy, options.SortOrder)).
		Find(&records)

	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

// Count returns the total number of visit rows.
func (r *Repository) Count() (int, error) {
	var count int64

	result := r.DB.Model(&Record{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

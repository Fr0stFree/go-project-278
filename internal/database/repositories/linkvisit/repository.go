// Package linkvisit provides services for managing link visits and interacting with the link visit repository.
package linkvisit

import (
	"fmt"
	"shortener/internal/database/postgres"
)

// Repository is a PostgreSQL implementation of AbstractRepository.
type Repository struct {
	*postgres.DataBase
}

// NewRepository creates a new instance of the Repository with the provided database connection.
func NewRepository(db *postgres.DataBase) *Repository {
	return &Repository{db}
}

// CreateOne creates a new link visit record in the database.
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

// GetMany retrieves multiple link visit records from the database based on the provided FilterOpts.
func (r *Repository) GetMany(options FilterOpts) ([]Record, error) {
	records := make([]Record, 0)

	result := r.DB.
		Limit(options.limit).
		Offset(options.offset).
		Order(fmt.Sprintf("%s %s", options.sortBy, options.sortOrder)).
		Find(&records)

	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

// Count returns the total number of link visit records in the database.
func (r *Repository) Count() (int, error) {
	var count int64

	result := r.DB.Model(&Record{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

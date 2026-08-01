package link

import (
	"errors"
	"fmt"
	"shortener/internal/database/postgres"

	"gorm.io/gorm"
)

// Repository is a PostgreSQL implementation of storage.AbstractLinkRepository.
type Repository struct {
	*postgres.DataBase
}

// NewRepository creates a new instance of the Repository with the provided database connection.
func NewRepository(db *postgres.DataBase) *Repository {
	return &Repository{db}
}

// CreateOne creates the original link and its corresponding short name in the database.
func (r *Repository) CreateOne(insert Insert) (Record, error) {
	link := Record{
		OriginalURL: insert.OriginalURL,
		ShortName:   insert.ShortName,
	}

	result := r.DB.Create(&link)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return Record{}, ErrShortNameAlreadyTaken
		}

		return Record{}, result.Error
	}

	return link, nil
}

// GetByID retrieves the original URL corresponding to the given link ID from the database.
func (r *Repository) GetByID(ID int) (Record, error) {
	var link Record

	result := r.DB.First(&link, ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Record{}, ErrNotFound
		}

		return Record{}, result.Error
	}

	return link, nil
}

// GetMany lists all links in the database with optional
// filtering, sorting, and pagination based on the provided FilterOpts.
func (r *Repository) GetMany(options FilterOpts) ([]Record, error) {
	records := make([]Record, 0)

	statement := r.DB.Model(&Record{})
	if len(options.shortNames) > 0 {
		statement = statement.Where("short_name IN ?", options.shortNames)
	}

	result := statement.
		Limit(options.limit).
		Offset(options.offset).
		Order(fmt.Sprintf("%s %s", options.sortBy, options.sortOrder)).
		Find(&records)

	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

// Count returns the total number of links in the database.
func (r *Repository) Count() (int, error) {
	var count int64

	result := r.DB.Model(&Record{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

// UpdateByID updates the original URL and short name for the given link ID in the database.
func (r *Repository) UpdateByID(ID int, update Update) (Record, error) {
	var link Record

	result := r.DB.First(&link, ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Record{}, ErrNotFound
		}

		return Record{}, result.Error
	}

	link.OriginalURL = update.OriginalURL
	link.ShortName = update.ShortName

	result = r.DB.Save(&link)
	if result.Error != nil {
		return Record{}, result.Error
	}

	return link, nil
}

// DeleteByID deletes the link with the given ID from the database.
func (r *Repository) DeleteByID(ID int) error {
	result := r.DB.Where("id = ?", ID).Delete(&Record{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

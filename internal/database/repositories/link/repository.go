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

func NewRepository (db *postgres.DataBase) *Repository {
	return &Repository{db}
}

// CreateOne creates the original link and its corresponding short name in PostgreSQL.
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

// GetByID retrieves the original URL corresponding to the given link ID from PostgreSQL.
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

// GetMany lists all links in the PostgreSQL database.
func (r *Repository) GetMany(options filterOpts) ([]Record, error) {
	links := make([]Record, 0)

	result := r.DB.
		Order(fmt.Sprintf("%s %s", options.sortBy, options.sortOrder)).
		Offset(options.offset).
		Limit(options.limit).
		Find(&links)

	if result.Error != nil {
		return nil, result.Error
	}

	return links, nil
}

// UpdateByID updates the original URL and short name for the given link ID in PostgreSQL.
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

// DeleteByID deletes the link with the given ID from PostgreSQL.
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

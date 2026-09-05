package link

import (
	"errors"
	"fmt"
	"shortener/internal/db"
	"shortener/internal/db/models"

	"gorm.io/gorm"
)

// Repository stores shortened links in PostgreSQL.
type Repository struct {
	*db.DataBase
}

// NewRepository creates a link repository backed by the provided database.
func NewRepository(database *db.DataBase) *Repository {
	return &Repository{database}
}

// CreateOne inserts a shortened link row.
func (r *Repository) CreateOne(insert Insert) (Record, error) {
	record := Record{
		OriginalURL: insert.OriginalURL,
		ShortName:   insert.ShortName,
	}

	result := r.DB.Create(&record)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return Record{}, models.ErrObjectAlreadyExists
		}

		return Record{}, result.Error
	}

	return record, nil
}

// GetByID returns a link row by ID.
func (r *Repository) GetByID(ID uint) (Record, error) {
	var record Record

	result := r.DB.First(&record, ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Record{}, models.ErrObjectDoesNotExist
		}

		return Record{}, result.Error
	}

	return record, nil
}

// GetMany returns link rows matching the provided list options.
func (r *Repository) GetMany(options ListOptions) ([]Record, error) {
	records := make([]Record, 0)

	statement := r.DB.Model(&Record{})
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
func (r *Repository) Count() (int, error) {
	var count int64

	result := r.DB.Model(&Record{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

// UpdateByID replaces URL fields for a link row by ID.
func (r *Repository) UpdateByID(ID uint, update Update) (Record, error) {
	var record Record

	result := r.DB.First(&record, ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Record{}, models.ErrObjectDoesNotExist
		}

		return Record{}, result.Error
	}

	record.OriginalURL = update.OriginalURL
	record.ShortName = update.ShortName

	result = r.DB.Save(&record)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return Record{}, models.ErrObjectAlreadyExists
		}

		return Record{}, result.Error
	}

	return record, nil
}

// DeleteByID deletes a link row by ID.
func (r *Repository) DeleteByID(ID uint) error {
	result := r.DB.Where("id = ?", ID).Delete(&Record{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrObjectDoesNotExist
	}

	return nil
}

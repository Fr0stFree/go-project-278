// Package postgres provides PostgreSQL link repository implementations.
package postgres

import (
	"errors"
	"fmt"
	"shortener/internal/config"
	"shortener/internal/storage"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// LinkRepository is a PostgreSQL implementation of storage.AbstractLinkRepository.
type LinkRepository struct {
	config *config.Storage
	db     *gorm.DB
}

// NewLinkRepository creates a new PostgreSQL link repository with the provided DSN.
// It establishes a connection to the PostgreSQL database and performs necessary migrations.
func NewLinkRepository(config *config.Storage) (*LinkRepository, error) {
	sslMode := "disable"
	if config.IsSSLEnabled {
		sslMode = "require"
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	err = db.AutoMigrate(&storage.LinkRecord{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	dbSQL, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	dbSQL.SetMaxOpenConns(config.MaxOpenConnections)
	dbSQL.SetMaxIdleConns(config.MaxIdleConnections)
	dbSQL.SetConnMaxLifetime(config.ConnectionMaxLifetime)

	return &LinkRepository{
		config: config,
		db:     db,
	}, nil
}

// SaveLink saves the original link and its corresponding short name in PostgreSQL.
func (r *LinkRepository) SaveLink(insert storage.LinkInsert) (storage.LinkRecord, error) {
	link := storage.LinkRecord{
		OriginalURL: insert.OriginalURL,
		ShortName:   insert.ShortName,
	}
	result := r.db.Create(&link)
	if result.Error != nil {
		return storage.LinkRecord{}, result.Error
	}

	return link, nil
}

// GetLink retrieves the original URL corresponding to the given link ID from PostgreSQL.
func (r *LinkRepository) GetLink(ID int) (storage.LinkRecord, error) {
	var link storage.LinkRecord
	result := r.db.First(&link, ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return storage.LinkRecord{}, storage.ErrLinkNotFound
		}
		return storage.LinkRecord{}, result.Error
	}
	return link, nil
}

// ListLinks lists all links in the PostgreSQL database.
func (r *LinkRepository) ListLinks() ([]storage.LinkRecord, error) {
	links := make([]storage.LinkRecord, 0)
	result := r.db.Order("id ASC").Find(&links)
	if result.Error != nil {
		return nil, result.Error
	}

	return links, nil
}

// UpdateLink updates the original URL and short name for the given link ID in PostgreSQL.
func (r *LinkRepository) UpdateLink(ID int, update storage.LinkUpdate) (storage.LinkRecord, error) {
	var link storage.LinkRecord
	result := r.db.First(&link, ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return storage.LinkRecord{}, storage.ErrLinkNotFound
		}
		return storage.LinkRecord{}, result.Error
	}

	link.OriginalURL = update.OriginalURL
	link.ShortName = update.ShortName

	result = r.db.Save(&link)
	if result.Error != nil {
		return storage.LinkRecord{}, result.Error
	}

	return link, nil
}

// DeleteLink deletes the link with the given ID from PostgreSQL.
func (r *LinkRepository) DeleteLink(ID int) error {
	result := r.db.Where("id = ?", ID).Delete(&storage.LinkRecord{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return storage.ErrLinkNotFound
	}
	return nil
}

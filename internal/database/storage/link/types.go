// Package link stores shortened link records.
package link

import (
	"shortener/internal/database/storage"

	"gorm.io/gorm"
)

// Record maps a shortened link row.
type Record struct {
	gorm.Model
	OriginalURL string `gorm:"column:original_url;not null"`
	ShortName   string `gorm:"column:short_name;not null;uniqueIndex"`
}

// TableName returns the database table for link records.
func (Record) TableName() string {
	return "shortened_links"
}

// Insert contains values for creating a link row.
type Insert struct {
	OriginalURL string
	ShortName   string
}

// Update contains values for replacing a link row.
type Update struct {
	OriginalURL string
	ShortName   string
}

// Filters contains link-specific query filters.
type Filters struct {
	ShortNames []string
}

// ListOptions combines pagination, sorting, and link filters.
type ListOptions struct {
	storage.ListOptions
	Filters
}

// AbstractRepository describes storage operations required for links.
type AbstractRepository interface {
	CreateOne(insert Insert) (Record, error)
	GetByID(ID uint) (Record, error)
	GetMany(options ListOptions) ([]Record, error)
	Count() (int, error)
	UpdateByID(ID uint, update Update) (Record, error)
	DeleteByID(ID uint) error
}

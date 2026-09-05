// Package link stores shortened link records.
package link

import (
	"shortener/internal/db"

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
	db.ListOptions
	Filters
}

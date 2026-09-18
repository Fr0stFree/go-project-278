// Package link stores shortened link records.
package link

import (
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

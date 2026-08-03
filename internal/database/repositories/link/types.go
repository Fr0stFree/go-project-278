// Package link provides types and interfaces for managing shortened links in the storage system.
package link

import "gorm.io/gorm"

// Record represents the output data returned when retrieving a link from the storage system.
type Record struct {
	gorm.Model
	OriginalURL string `gorm:"column:original_url;not null"`
	ShortName   string `gorm:"column:short_name;not null;uniqueIndex"`
}

// TableName specifies the table name for the Record struct in the database.
func (Record) TableName() string {
	return "shortened_links"
}

// Insert represents the input data required to save a link in the storage system.
type Insert struct {
	OriginalURL string
	ShortName   string
}

// Update represents the input data required to update a link in the storage system.
type Update struct {
	OriginalURL string
	ShortName   string
}

// AbstractRepository defines the interface for link storage systems.
type AbstractRepository interface {
	CreateOne(insert Insert) (Record, error)
	GetByID(ID uint) (Record, error)
	GetMany(options FilterOpts) ([]Record, error)
	Count() (int, error)
	UpdateByID(ID uint, update Update) (Record, error)
	DeleteByID(ID uint) error
}

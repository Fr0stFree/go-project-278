package linkvisit

import (
	"shortener/internal/database/storage"
	"shortener/internal/database/storage/link"

	"gorm.io/gorm"
)

// Record maps a redirect visit row.
type Record struct {
	gorm.Model
	LinkID    uint        `gorm:"column:link_id;not null"`
	Link      link.Record `gorm:"foreignKey:LinkID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IP        string      `gorm:"column:ip;not null"`
	UserAgent string      `gorm:"column:user_agent;not null"`
	Status    uint        `gorm:"column:status;not null"`
	Referrer  string      `gorm:"column:referrer"`
}

// TableName returns the database table for visit records.
func (Record) TableName() string {
	return "shortened_link_visits"
}

// Insert contains values for creating a visit row.
type Insert struct {
	LinkID    uint
	IP        string
	UserAgent string
	Status    uint
	Referrer  string
}

// Filters contains visit-specific query filters.
type Filters struct {
	LinkIDs []uint
}

// ListOptions combines pagination, sorting, and visit filters.
type ListOptions struct {
	storage.ListOptions
	Filters
}

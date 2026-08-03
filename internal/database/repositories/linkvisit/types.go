package linkvisit

import (
	"shortener/internal/database/repositories/link"

	"gorm.io/gorm"
)

// Record represents a record of a link visit, including the associated link, IP address, user agent, and status.
type Record struct {
	gorm.Model
	LinkID    uint        `gorm:"column:link_id;not null"`
	Link      link.Record `gorm:"foreignKey:LinkID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IP        string      `gorm:"column:ip;not null"`
	UserAgent string      `gorm:"column:user_agent;not null"`
	Status    uint        `gorm:"column:status;not null"`
	Referrer  string      `gorm:"column:referrer"`
}

// TableName specifies the table name for the Record struct in the database.
func (Record) TableName() string {
	return "shortened_link_visits"
}

// Insert represents the input data required to save a link visit in the storage system.
type Insert struct {
	LinkID    uint
	IP        string
	UserAgent string
	Status    uint
	Referrer  string
}

// AbstractRepository defines the interface for link visit storage systems.
type AbstractRepository interface {
	CreateOne(insert Insert) (Record, error)
	GetMany(options FilterOpts) ([]Record, error)
	Count() (int, error)
}

// Package storage defines the interface, errors and common types for link storage systems.
package storage

// LinkRecord represents the output data returned when retrieving a link from the storage system.
type LinkRecord struct {
	ID          int    `gorm:"primaryKey"`
	OriginalURL string `gorm:"column:original_url;not null"`
	ShortName   string `gorm:"column:short_name;not null"`
}

func (LinkRecord) TableName() string {
	return "shortened_links"
}

// LinkInsert represents the input data required to save a link in the storage system.
type LinkInsert struct {
	OriginalURL string
	ShortName   string
}

type LinkUpdate struct {
	OriginalURL string
	ShortName   string
}

// AbstractLinkRepository defines the interface for link storage systems.
type AbstractLinkRepository interface {
	SaveLink(insert LinkInsert) (LinkRecord, error)
	GetLink(ID int) (LinkRecord, error)
	ListLinks() ([]LinkRecord, error)
	UpdateLink(ID int, update LinkUpdate) (LinkRecord, error)
	DeleteLink(ID int) error
}

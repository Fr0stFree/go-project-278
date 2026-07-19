// Package storage defines the interface, errors and common types for link storage systems.
package storage

// LinkDBOut represents the output data returned when retrieving a link from the storage system.
type LinkDBOut struct {
	ID          int
	OriginalURL string
	ShortName   string
}

// LinkDBIn represents the input data required to save a link in the storage system.
type LinkDBIn struct {
	OriginalURL string
	ShortName   string
}

// Storage defines the interface for a link storage system.
type Storage interface {
	SaveLink(link LinkDBIn) (LinkDBOut, error)
	GetLink(ID int) (LinkDBOut, error)
	ListLinks() ([]LinkDBOut, error)
	UpdateLink(ID int, update LinkDBIn) (LinkDBOut, error)
}

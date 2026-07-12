// Package memory provides an in-memory implementation of the Storage interface.
package memory

import (
	"shortener/internal/storage"
)

// Storage is an in-memory implementation of the Storage interface.
type Storage struct {
	data       map[int]storage.LinkDBOut
	sequenceID int
}

// NewStorage creates a new instance of the in-memory storage.
func NewStorage() *Storage {
	return &Storage{
		data:       make(map[int]storage.LinkDBOut),
		sequenceID: 0,
	}
}

// SaveLink saves the original link and its corresponding short name in memory.
func (s *Storage) SaveLink(link storage.LinkDBIn) (storage.LinkDBOut, error) {
	s.sequenceID++
	obj := storage.LinkDBOut{
		ID:          s.sequenceID,
		OriginalURL: link.OriginalURL,
		ShortName:   link.ShortName,
	}
	s.data[s.sequenceID] = obj

	return obj, nil
}

// GetLink retrieves the original URL corresponding to the given link ID from memory.
func (s *Storage) GetLink(ID int) (storage.LinkDBOut, error) {
	link, exists := s.data[ID]
	if !exists {
		return storage.LinkDBOut{}, storage.ErrLinkNotFound
	}

	return link, nil
}

// ListLinks retrieves a list of all shortened links stored in memory.
func (s *Storage) ListLinks() ([]storage.LinkDBOut, error) {
	links := make([]storage.LinkDBOut, 0, len(s.data))
	for _, link := range s.data {
		links = append(links, link)
	}

	return links, nil
}

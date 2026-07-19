// Package memory provides in-memory link repository implementations.
package memory

import (
	"shortener/internal/storage"
)

// LinkRepository is an in-memory implementation of storage.AbstractLinkRepository.
type LinkRepository struct {
	data       map[int]storage.LinkDBOut
	sequenceID int
}

// NewLinkRepository creates a new in-memory link repository.
func NewLinkRepository() *LinkRepository {
	return &LinkRepository{
		data:       make(map[int]storage.LinkDBOut),
		sequenceID: 0,
	}
}

// SaveLink saves the original link and its corresponding short name in memory.
func (r *LinkRepository) SaveLink(link storage.LinkDBIn) (storage.LinkDBOut, error) {
	r.sequenceID++
	obj := storage.LinkDBOut{
		ID:          r.sequenceID,
		OriginalURL: link.OriginalURL,
		ShortName:   link.ShortName,
	}
	r.data[r.sequenceID] = obj

	return obj, nil
}

// GetLink retrieves the original URL corresponding to the given link ID from memory.
func (r *LinkRepository) GetLink(ID int) (storage.LinkDBOut, error) {
	link, exists := r.data[ID]
	if !exists {
		return storage.LinkDBOut{}, storage.ErrLinkNotFound
	}

	return link, nil
}

// ListLinks retrieves a list of all shortened links stored in memory.
func (r *LinkRepository) ListLinks() ([]storage.LinkDBOut, error) {
	links := make([]storage.LinkDBOut, 0, len(r.data))
	for _, link := range r.data {
		links = append(links, link)
	}

	return links, nil
}

func (r *LinkRepository) UpdateLink(ID int, update storage.LinkDBIn) (storage.LinkDBOut, error) {
	link, exists := r.data[ID]
	if !exists {
		return storage.LinkDBOut{}, storage.ErrLinkNotFound
	}
	link.OriginalURL = update.OriginalURL
	link.ShortName = update.ShortName
	r.data[ID] = link

	return link, nil
}

func (r *LinkRepository) DeleteLink(ID int) error {
	_, exists := r.data[ID]
	if !exists {
		return storage.ErrLinkNotFound
	}
	delete(r.data, ID)

	return nil
}

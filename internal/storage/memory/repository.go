// Package memory provides in-memory repository implementations.
package memory

import (
	"shortener/internal/storage"
)

// Repository is an in-memory implementation of storage.LinkRepository.
type Repository struct {
	data       map[int]storage.LinkDBOut
	sequenceID int
}

// NewRepository creates a new in-memory link repository.
func NewRepository() *Repository {
	return &Repository{
		data:       make(map[int]storage.LinkDBOut),
		sequenceID: 0,
	}
}

// SaveLink saves the original link and its corresponding short name in memory.
func (r *Repository) SaveLink(link storage.LinkDBIn) (storage.LinkDBOut, error) {
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
func (r *Repository) GetLink(ID int) (storage.LinkDBOut, error) {
	link, exists := r.data[ID]
	if !exists {
		return storage.LinkDBOut{}, storage.ErrLinkNotFound
	}

	return link, nil
}

// ListLinks retrieves a list of all shortened links stored in memory.
func (r *Repository) ListLinks() ([]storage.LinkDBOut, error) {
	links := make([]storage.LinkDBOut, 0, len(r.data))
	for _, link := range r.data {
		links = append(links, link)
	}

	return links, nil
}

func (r *Repository) UpdateLink(ID int, update storage.LinkDBIn) (storage.LinkDBOut, error) {
	link, exists := r.data[ID]
	if !exists {
		return storage.LinkDBOut{}, storage.ErrLinkNotFound
	}
	link.OriginalURL = update.OriginalURL
	link.ShortName = update.ShortName
	r.data[ID] = link

	return link, nil
}

func (r *Repository) DeleteLink(ID int) error {
	_, exists := r.data[ID]
	if !exists {
		return storage.ErrLinkNotFound
	}
	delete(r.data, ID)

	return nil
}

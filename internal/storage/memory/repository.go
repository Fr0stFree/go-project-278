// Package memory provides in-memory link repository implementations.
package memory

import (
	"shortener/internal/storage"
)

// LinkRepository is an in-memory implementation of storage.AbstractLinkRepository.
type LinkRepository struct {
	data       map[int]storage.LinkRecord
	sequenceID int
}

// NewLinkRepository creates a new in-memory link repository.
func NewLinkRepository() *LinkRepository {
	return &LinkRepository{
		data:       make(map[int]storage.LinkRecord),
		sequenceID: 0,
	}
}

// SaveLink saves the original link and its corresponding short name in memory.
func (r *LinkRepository) SaveLink(insert storage.LinkInsert) (storage.LinkRecord, error) {
	r.sequenceID++
	link := storage.LinkRecord{
		ID:          r.sequenceID,
		OriginalURL: insert.OriginalURL,
		ShortName:   insert.ShortName,
	}
	r.data[r.sequenceID] = link

	return link, nil
}

// GetLink retrieves the original URL corresponding to the given link ID from memory.
func (r *LinkRepository) GetLink(ID int) (storage.LinkRecord, error) {
	link, exists := r.data[ID]
	if !exists {
		return storage.LinkRecord{}, storage.ErrLinkNotFound
	}

	return link, nil
}

// ListLinks retrieves a list of all shortened links stored in memory.
func (r *LinkRepository) ListLinks() ([]storage.LinkRecord, error) {
	links := make([]storage.LinkRecord, 0, len(r.data))
	for _, link := range r.data {
		links = append(links, link)
	}

	return links, nil
}

func (r *LinkRepository) UpdateLink(ID int, update storage.LinkUpdate) (storage.LinkRecord, error) {
	link, exists := r.data[ID]
	if !exists {
		return storage.LinkRecord{}, storage.ErrLinkNotFound
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

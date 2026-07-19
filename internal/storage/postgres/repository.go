// Package postgres provides PostgreSQL link repository implementations.
package postgres

import (
	"shortener/internal/config"
	"shortener/internal/storage"
)

// LinkRepository is a PostgreSQL implementation of storage.AbstractLinkRepository.
type LinkRepository struct {
	config *config.Storage
}

// NewLinkRepository creates a new PostgreSQL link repository with the provided DSN.
func NewLinkRepository(config *config.Storage) *LinkRepository {
	return &LinkRepository{
		config: config,
	}
}

func (r *LinkRepository) SaveLink(link storage.LinkDBIn) (storage.LinkDBOut, error) {
	return storage.LinkDBOut{}, nil
}

func (r *LinkRepository) GetLink(ID int) (storage.LinkDBOut, error) {
	return storage.LinkDBOut{}, nil
}

func (r *LinkRepository) ListLinks() ([]storage.LinkDBOut, error) {
	return nil, nil
}

func (r *LinkRepository) UpdateLink(ID int, update storage.LinkDBIn) (storage.LinkDBOut, error) {
	return storage.LinkDBOut{}, nil
}

func (r *LinkRepository) DeleteLink(ID int) error {
	return nil
}

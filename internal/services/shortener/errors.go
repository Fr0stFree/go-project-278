package shortener

import (
	"errors"
	"shortener/internal/database/storage"
)

var (
	// ErrLinkNotFound is returned when a requested link is not found in the storage.
	ErrLinkNotFound = errors.New("link not found")
	// ErrShortNameAlreadyTaken is returned when attempting to create a link with a short name that is already in use.
	ErrShortNameAlreadyTaken = errors.New("shortname already taken")
)

func mapStorageErrorToServiceError(err error) error {
	switch {
	case errors.Is(err, storage.ErrObjectDoesNotExist):
		return ErrLinkNotFound
	case errors.Is(err, storage.ErrObjectAlreadyExists):
		return ErrShortNameAlreadyTaken
	default:
		return err
	}
}

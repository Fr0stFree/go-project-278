package storage

import "errors"

var (
	// ErrLinkNotFound is returned when a requested link is not found in the storage.
	ErrLinkNotFound = errors.New("link not found")
)

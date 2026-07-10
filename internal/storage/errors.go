package storage

import "errors"

var (
	// ErrURLNotFound is returned when a requested URL is not found in the storage.
	ErrURLNotFound = errors.New("URL not found")
)

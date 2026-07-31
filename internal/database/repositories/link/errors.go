package link

import "errors"

var (
	// ErrNotFound is returned when a requested link is not found in the storage.
	ErrNotFound = errors.New("link not found")

	ErrShortNameAlreadyTaken = errors.New("shortname already taken")
)

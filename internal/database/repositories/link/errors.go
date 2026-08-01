package link

import "errors"

var (
	// ErrNotFound is returned when a requested link is not found in the storage.
	ErrNotFound = errors.New("link not found")
	// ErrShortNameAlreadyTaken is returned when attempting to create a link with a short name that is already in use.
	ErrShortNameAlreadyTaken = errors.New("shortname already taken")
)

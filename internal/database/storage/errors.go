// Package storage defines common repository options and errors.
package storage

import "errors"

var (
	// ErrObjectDoesNotExist reports that a requested row was not found.
	ErrObjectDoesNotExist = errors.New("object does not exist")
	// ErrObjectAlreadyExists reports that a unique row already exists.
	ErrObjectAlreadyExists = errors.New("object already exists")
)

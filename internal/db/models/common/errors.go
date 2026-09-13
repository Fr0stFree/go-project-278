// Package common defines contracts shared by persistence implementations.
package common

import "errors"

var (
	// ErrObjectDoesNotExist reports that a requested record was not found.
	ErrObjectDoesNotExist = errors.New("object does not exist")
	// ErrObjectAlreadyExists reports that a unique record already exists.
	ErrObjectAlreadyExists = errors.New("object already exists")
)

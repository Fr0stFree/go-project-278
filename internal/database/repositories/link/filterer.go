package link

import (
	"fmt"
)

const (
	defaultLimit     = 10
	defaultOffset    = 0
	defaultSortBy    = "id"
	defaultSortOrder = "ASC"
)

var sortOptions = map[string]string{
	"id":           "id",
	"original_url": "original_url",
	"short_name":   "short_name",
	"short_url":    "short_name",
	"created_at":   "created_at",
}

var sortOrderOptions = map[string]string{
	"ASC":  "ASC",
	"DESC": "DESC",
}

// FilterOpts represents the filtering options for retrieving links from the database.
type FilterOpts struct {
	limit     int
	offset    int
	sortBy    string
	sortOrder string
}

// NewFilterOpts creates a new instance of FilterOpts with default values.
func NewFilterOpts() *FilterOpts {
	return &FilterOpts{
		limit:     defaultLimit,
		offset:    defaultOffset,
		sortBy:    defaultSortBy,
		sortOrder: defaultSortOrder,
	}
}

// Range returns the range of records to retrieve based on the current limit and offset values.
func (o *FilterOpts) Range() (from, to int) {
	from = o.offset
	to = o.offset + o.limit - 1

	return from, to
}

// WithRange sets the range of records to retrieve based on the provided from and to indices.
func (o *FilterOpts) WithRange(from, to int) error {
	if from < 0 || to < 0 || from > to {
		return fmt.Errorf("invalid range: from=%d, to=%d", from, to)
	}

	o.limit = to - from + 1
	o.offset = from

	return nil
}

// WithSort sets the sorting criteria for the records to retrieve based on the provided sortBy and sortOrder values.
func (o *FilterOpts) WithSort(sortBy, sortOrder string) error {
	sortOption, exists := sortOptions[sortBy]
	if !exists {
		return fmt.Errorf("invalid sort by: %s", sortBy)
	}

	orderOption, exists := sortOrderOptions[sortOrder]
	if !exists {
		return fmt.Errorf("invalid sort order: %s", sortOrder)
	}

	o.sortBy = sortOption
	o.sortOrder = orderOption

	return nil
}

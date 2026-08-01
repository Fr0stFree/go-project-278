package linkvisit

import (
	"fmt"
)

const (
	defaultLimit     = 10
	defaultOffset    = 0
	defaultSortBy    = "created_at"
	defaultSortOrder = "DESC"
)

// FilterOpts represents the filtering options for retrieving link visits from the database.
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
func (o *FilterOpts) WithRange(from, to int) (*FilterOpts, error) {
	if from < 0 || to < 0 || from > to {
		return nil, fmt.Errorf("invalid range: from=%d, to=%d", from, to)
	}

	o.offset = from
	o.limit = to - from + 1

	return o, nil
}

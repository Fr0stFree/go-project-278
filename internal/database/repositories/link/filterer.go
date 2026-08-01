package link

const (
	defaultLimit     = 10
	defaultOffset    = 0
	defaultSortBy    = "id"
	defaultSortOrder = "ASC"
)

// FilterOpts represents the filtering options for retrieving links from the database.
type FilterOpts struct {
	limit     int
	offset    int
	sortBy    string
	sortOrder string
}

// NewFilterOpts creates a new instance of FilterOpts with default values.
func NewFilterOpts() FilterOpts {
	return FilterOpts{
		limit:     defaultLimit,
		offset:    defaultOffset,
		sortBy:    defaultSortBy,
		sortOrder: defaultSortOrder,
	}
}

// WithRange sets the range of records to retrieve based on the provided from and to indices.
func (o *FilterOpts) WithRange(from, to int) {
	o.limit = to - from + 1
	o.offset = from
}

// WithSort sets the sorting criteria for the records to retrieve based on the provided sortBy and sortOrder values.
func (o *FilterOpts) WithSort(sortBy, sortOrder string) {
	o.sortBy = sortBy
	o.sortOrder = sortOrder
}

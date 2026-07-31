package link

const (
	defaultLimit     = 10
	defaultOffset    = 0
	defaultSortBy    = "id"
	defaultSortOrder = "ASC"
)

type filterOpts struct {
	limit     int
	offset    int
	sortBy    string
	sortOrder string
}

func NewFilterOpts() filterOpts {
	return filterOpts{
		limit:     defaultLimit,
		offset:    defaultOffset,
		sortBy:    defaultSortBy,
		sortOrder: defaultSortOrder,
	}
}

func (o *filterOpts) WithRange(from, to int) {
	o.limit = to - from + 1
	o.offset = from
}

func (o *filterOpts) WithSort(sortBy, sortOrder string) {
	o.sortBy = sortBy
	o.sortOrder = sortOrder
}

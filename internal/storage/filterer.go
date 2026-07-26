package storage

const (
	defaultLimit     = 10
	defaultOffset    = 0
	defaultSortBy    = "id"
	defaultSortOrder = "ASC"
)

type listLinksOptions struct {
	limit     int
	offset    int
	sortBy    string
	sortOrder string
}

func NewListLinksOptions() listLinksOptions {
	return listLinksOptions{
		limit:     defaultLimit,
		offset:    defaultOffset,
		sortBy:    defaultSortBy,
		sortOrder: defaultSortOrder,
	}
}

func (o *listLinksOptions) WithRange(from, to int) {
	o.limit = to - from + 1
	o.offset = from
}

func (o *listLinksOptions) WithSort(sortBy, sortOrder string) {
	o.sortBy = sortBy
	o.sortOrder = sortOrder
}

package shortener

import (
	"fmt"
	"strings"

	"shortener/internal/database/storage"
	"shortener/internal/database/storage/link"
	"shortener/internal/database/storage/linkvisit"
)

// ListOptionsBuilder validates pagination and sorting options.
type ListOptionsBuilder struct {
	sortFields map[string]string
	options    storage.ListOptions
	err        error
}

// Error returns the first validation error.
func (b *ListOptionsBuilder) Error() error {
	return b.err
}

// WithRange sets an inclusive zero-based result range.
func (b *ListOptionsBuilder) WithRange(from, to int) *ListOptionsBuilder {
	if b.err != nil {
		return b
	}

	if from < 0 {
		b.err = fmt.Errorf("range start must be non-negative: %d", from)

		return b
	}

	if to < from {
		b.err = fmt.Errorf(
			"range end must be greater than or equal to range start: %d < %d",
			to,
			from,
		)

		return b
	}

	b.options.Offset = from
	b.options.Limit = to - from + 1

	return b
}

// WithSort sets a whitelisted sort field and direction.
func (b *ListOptionsBuilder) WithSort(field, order string) *ListOptionsBuilder {
	if b.err != nil {
		return b
	}

	sortBy, ok := b.sortFields[field]
	if !ok {
		b.err = fmt.Errorf("unsupported sort field: %q", field)

		return b
	}

	order = strings.ToUpper(order)
	switch order {
	case "ASC", "DESC":
	default:
		b.err = fmt.Errorf("unsupported sort order: %q", order)

		return b
	}

	b.options.SortBy = sortBy
	b.options.SortOrder = order

	return b
}

// Range returns the current inclusive result range.
func (b *ListOptionsBuilder) Range() (int, int) {
	return b.options.Offset, b.options.Offset + b.options.Limit - 1
}

// LinkListOptionsBuilder builds options for link list queries.
type LinkListOptionsBuilder struct {
	*ListOptionsBuilder
	filters link.Filters
}

// NewLinkListOptionsBuilder creates link list options with defaults.
func NewLinkListOptionsBuilder() *LinkListOptionsBuilder {
	var sortFields = map[string]string{
		"id":           "id",
		"original_url": "original_url",
		"short_name":   "short_name",
		"short_url":    "short_name",
		"created_at":   "created_at",
	}

	defaultSortBy := "id"
	defaultSortOrder := "DESC"
	defaultLimit := 10
	defaultOffset := 0

	return &LinkListOptionsBuilder{
		ListOptionsBuilder: &ListOptionsBuilder{
			sortFields: sortFields,
			options: storage.ListOptions{
				Limit:     defaultLimit,
				Offset:    defaultOffset,
				SortBy:    defaultSortBy,
				SortOrder: defaultSortOrder,
			},
		},
	}
}

// WithShortNames filters links by short name.
func (b *LinkListOptionsBuilder) WithShortNames(shortNames ...string) *LinkListOptionsBuilder {
	if b.Error() != nil {
		return b
	}

	b.filters.ShortNames = append(b.filters.ShortNames, shortNames...)

	return b
}

func (b *LinkListOptionsBuilder) build() link.ListOptions {
	return link.ListOptions{
		ListOptions: b.options,
		Filters:     b.filters,
	}
}

// LinkVisitListOptionsBuilder builds options for visit list queries.
type LinkVisitListOptionsBuilder struct {
	*ListOptionsBuilder
	filters linkvisit.Filters
}

// NewLinkVisitListOptionsBuilder creates visit list options with defaults.
func NewLinkVisitListOptionsBuilder() *LinkVisitListOptionsBuilder {
	var sortFields = map[string]string{
		"id":         "id",
		"link_id":    "link_id",
		"created_at": "created_at",
	}

	const (
		defaultSortBy    = "created_at"
		defaultSortOrder = "DESC"
		defaultLimit     = 10
		defaultOffset    = 0
	)

	return &LinkVisitListOptionsBuilder{
		ListOptionsBuilder: &ListOptionsBuilder{
			sortFields: sortFields,
			options: storage.ListOptions{
				Limit:     defaultLimit,
				Offset:    defaultOffset,
				SortBy:    defaultSortBy,
				SortOrder: defaultSortOrder,
			},
		},
	}
}

// WithLinkIDs filters visits by link ID.
func (b *LinkVisitListOptionsBuilder) WithLinkIDs(linkIDs ...uint) *LinkVisitListOptionsBuilder {
	if b.Error() != nil {
		return b
	}

	b.filters.LinkIDs = append(b.filters.LinkIDs, linkIDs...)

	return b
}

func (b *LinkVisitListOptionsBuilder) build() linkvisit.ListOptions {
	return linkvisit.ListOptions{
		ListOptions: b.options,
		Filters:     b.filters,
	}
}

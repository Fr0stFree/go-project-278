package shortener

import (
	"fmt"
	"strings"

	"shortener/internal/database/storage"
	"shortener/internal/database/storage/link"
	"shortener/internal/database/storage/linkvisit"
)

// ListOptionsBuilder validates and builds common list options.
type ListOptionsBuilder struct {
	sortFields map[string]string
	options    storage.ListOptions
	err        error
}

// Error returns the first validation error encountered by the builder.
func (b *ListOptionsBuilder) Error() error {
	return b.err
}

// WithRange sets the inclusive range of records to retrieve.
func (b *ListOptionsBuilder) WithRange(from, to int) *ListOptionsBuilder {
	if b.err != nil {
		return b
	}

	if from < 0 {
		b.err = NewValidationError(fmt.Sprintf("range start must be non-negative: %d", from), "range")

		return b
	}

	if to < from {
		b.err = NewValidationError(fmt.Sprintf("range end must be greater than or equal to range start: %d < %d", to, from), "range")

		return b
	}

	b.options.Offset = from
	b.options.Limit = to - from + 1

	return b
}

// WithSort validates and sets the sort field and order.
func (b *ListOptionsBuilder) WithSort(field, order string) *ListOptionsBuilder {
	if b.err != nil {
		return b
	}

	sortBy, ok := b.sortFields[field]
	if !ok {
		b.err = NewValidationError(fmt.Sprintf("unsupported sort field: %q", field), "sort")

		return b
	}

	order = strings.ToUpper(order)
	switch order {
	case "ASC", "DESC":
	default:
		b.err = NewValidationError(fmt.Sprintf("unsupported sort order: %q", order), "sort")

		return b
	}

	b.options.SortBy = sortBy
	b.options.SortOrder = order

	return b
}

// Range returns the inclusive range of records.
func (b *ListOptionsBuilder) Range() (int, int) {
	return b.options.Offset, b.options.Offset + b.options.Limit - 1
}

// LinkListOptionsBuilder builds options for listing links.
type LinkListOptionsBuilder struct {
	*ListOptionsBuilder
	filters link.Filters
}

// NewLinkListOptionsBuilder creates a new builder for link list options.
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

// WithShortNames sets the short names to filter by.
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

// LinkVisitListOptionsBuilder builds options for listing link visits.
type LinkVisitListOptionsBuilder struct {
	*ListOptionsBuilder
	filters linkvisit.Filters
}

// NewLinkVisitListOptionsBuilder creates a new builder for link visit list options.
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

// WithLinkIDs sets the link IDs to filter by.
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

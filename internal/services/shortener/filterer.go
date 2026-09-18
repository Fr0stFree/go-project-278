package shortener

import (
	"fmt"
	"strings"
)

// SortDirection describes the order of a list result.
type SortDirection string

// Supported sort directions.
const (
	SortAscending  SortDirection = "ASC"
	SortDescending SortDirection = "DESC"
)

// LinkSortField is a persistence-independent link sort key.
type LinkSortField string

// Supported link sort fields.
const (
	LinkSortByID          LinkSortField = "id"
	LinkSortByOriginalURL LinkSortField = "original_url"
	LinkSortByShortName   LinkSortField = "short_name"
	LinkSortByCreatedAt   LinkSortField = "created_at"
)

// LinkVisitSortField is a persistence-independent visit sort key.
type LinkVisitSortField string

// Supported visit sort fields.
const (
	LinkVisitSortByID        LinkVisitSortField = "id"
	LinkVisitSortByLinkID    LinkVisitSortField = "link_id"
	LinkVisitSortByCreatedAt LinkVisitSortField = "created_at"
)

// ListOptions contains pagination and ordering shared by list operations.
type ListOptions struct {
	Limit     int
	Offset    int
	SortOrder SortDirection
}

// LinkListOptions contains link query criteria.
type LinkListOptions struct {
	ListOptions
	SortBy     LinkSortField
	ShortNames []string
}

// LinkVisitListOptions contains visit query criteria.
type LinkVisitListOptions struct {
	ListOptions
	SortBy  LinkVisitSortField
	LinkIDs []uint
}

// ListOptionsBuilder validates common list options.
type ListOptionsBuilder struct {
	sortFields map[string]struct{}
	sortBy     string
	options    ListOptions
	maxLimit   int
	err        error
}

// WithRange sets the inclusive range of records to retrieve.
func (b *ListOptionsBuilder) WithRange(from, to int) {
	if b.err != nil {
		return
	}

	if from < 0 {
		b.err = NewValidationError(fmt.Sprintf("range start must be non-negative: %d", from), "range")

		return
	}

	if to < from {
		b.err = NewValidationError(fmt.Sprintf("range end must be greater than or equal to range start: %d < %d", to, from), "range")

		return
	}

	limit := to - from + 1
	if limit > b.maxLimit {
		b.err = NewValidationError(fmt.Sprintf("range must contain at most %d records, got %d", b.maxLimit, limit), "range")

		return
	}

	b.options.Offset = from
	b.options.Limit = limit
}

// WithSort validates and sets the sort field and order.
func (b *ListOptionsBuilder) WithSort(field, order string) {
	if b.err != nil {
		return
	}

	_, ok := b.sortFields[field]
	if !ok {
		b.err = NewValidationError(fmt.Sprintf("unsupported sort field: %q", field), "sort")

		return
	}

	order = strings.ToUpper(order)
	if order != string(SortAscending) && order != string(SortDescending) {
		b.err = NewValidationError(fmt.Sprintf("unsupported sort order: %q", order), "sort")

		return
	}

	b.sortBy = field
	b.options.SortOrder = SortDirection(order)
}

// Range returns the inclusive range of records.
func (b *ListOptionsBuilder) Range() (int, int) {
	return b.options.Offset, b.options.Offset + b.options.Limit - 1
}

// LinkListOptionsBuilder builds options for listing links.
type LinkListOptionsBuilder struct {
	*ListOptionsBuilder
	shortNames []string
}

// NewLinkListOptionsBuilder creates a link options builder with safe defaults.
func NewLinkListOptionsBuilder() *LinkListOptionsBuilder {
	return &LinkListOptionsBuilder{ListOptionsBuilder: &ListOptionsBuilder{
		sortFields: map[string]struct{}{
			string(LinkSortByID):          {},
			string(LinkSortByOriginalURL): {},
			string(LinkSortByShortName):   {},
			string(LinkSortByCreatedAt):   {},
		},
		sortBy: string(LinkSortByID), maxLimit: 100,
		options: ListOptions{Limit: 10, SortOrder: SortDescending},
	}}
}

// WithShortNames sets the short names to filter by.
func (b *LinkListOptionsBuilder) WithShortNames(shortNames ...string) {
	if b.err == nil {
		b.shortNames = append(b.shortNames, shortNames...)
	}
}

func (b *LinkListOptionsBuilder) build() LinkListOptions {
	return LinkListOptions{ListOptions: b.options, SortBy: LinkSortField(b.sortBy), ShortNames: b.shortNames}
}

// LinkVisitListOptionsBuilder builds options for listing link visits.
type LinkVisitListOptionsBuilder struct {
	*ListOptionsBuilder
	linkIDs []uint
}

// NewLinkVisitListOptionsBuilder creates a visit options builder with safe defaults.
func NewLinkVisitListOptionsBuilder() *LinkVisitListOptionsBuilder {
	return &LinkVisitListOptionsBuilder{ListOptionsBuilder: &ListOptionsBuilder{
		sortFields: map[string]struct{}{
			string(LinkVisitSortByID):        {},
			string(LinkVisitSortByLinkID):    {},
			string(LinkVisitSortByCreatedAt): {},
		},
		sortBy: string(LinkVisitSortByCreatedAt), maxLimit: 100,
		options: ListOptions{Limit: 10, SortOrder: SortDescending},
	}}
}

// WithLinkIDs sets the link IDs to filter by.
func (b *LinkVisitListOptionsBuilder) WithLinkIDs(linkIDs ...uint) {
	if b.err == nil {
		b.linkIDs = append(b.linkIDs, linkIDs...)
	}
}

func (b *LinkVisitListOptionsBuilder) build() LinkVisitListOptions {
	return LinkVisitListOptions{ListOptions: b.options, SortBy: LinkVisitSortField(b.sortBy), LinkIDs: b.linkIDs}
}

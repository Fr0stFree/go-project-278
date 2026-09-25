package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkVisitListOptionsBuilderWithSort(t *testing.T) {
	tests := []struct {
		field     LinkVisitSortField
		direction string
		expected  SortDirection
	}{
		{field: LinkVisitSortByID, direction: "asc", expected: SortAscending},
		{field: LinkVisitSortByLinkID, direction: "DESC", expected: SortDescending},
		{field: LinkVisitSortByCreatedAt, direction: "asc", expected: SortAscending},
		{field: LinkVisitSortByIP, direction: "desc", expected: SortDescending},
		{field: LinkVisitSortByUserAgent, direction: "ASC", expected: SortAscending},
		{field: LinkVisitSortByStatus, direction: "desc", expected: SortDescending},
	}

	for _, test := range tests {
		t.Run(string(test.field), func(t *testing.T) {
			builder := NewLinkVisitListOptionsBuilder()
			builder.WithSort(string(test.field), test.direction)

			options := builder.build()

			require.NoError(t, builder.err)
			assert.Equal(t, test.field, options.SortBy)
			assert.Equal(t, test.expected, options.SortOrder)
		})
	}
}

func TestLinkVisitListOptionsBuilderRejectsInvalidSort(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		direction string
	}{
		{name: "unknown field", field: "referrer", direction: "ASC"},
		{name: "unknown direction", field: "status", direction: "sideways"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := NewLinkVisitListOptionsBuilder()
			builder.WithSort(test.field, test.direction)

			var validationErr *ValidationError
			require.ErrorAs(t, builder.err, &validationErr)
			assert.Equal(t, "sort", validationErr.Field)
		})
	}
}

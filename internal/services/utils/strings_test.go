package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSnakeCase(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected string
	}

	tests := []testCase{
		{
			name:     "should convert camel case",
			input:    "ShortName",
			expected: "short_name",
		},
		{
			name:     "should handle acronym at the end",
			input:    "OriginalURL",
			expected: "original_url",
		},
		{
			name:     "should handle acronym at the beginning",
			input:    "HTTPStatus",
			expected: "http_status",
		},
		{
			name:     "should handle ID suffix",
			input:    "LinkID",
			expected: "link_id",
		},
		{
			name:     "should handle single acronym",
			input:    "URL",
			expected: "url",
		},
		{
			name:     "should handle single word",
			input:    "Name",
			expected: "name",
		},
		{
			name:     "should preserve lowercase string",
			input:    "name",
			expected: "name",
		},
		{
			name:     "should preserve snake case",
			input:    "short_name",
			expected: "short_name",
		},
		{
			name:     "should handle empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ToSnakeCase(tt.input))
		})
	}
}

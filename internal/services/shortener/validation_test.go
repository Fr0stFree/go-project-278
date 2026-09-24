package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_ValidateOriginalURL(t *testing.T) {
	service := NewService(nil, nil)

	type subTest struct {
		name        string
		originalURL string
		expectedErr string
	}

	subTests := []subTest{
		{
			name:        "valid http URL",
			originalURL: "http://example.com",
			expectedErr: "",
		},
		{
			name:        "valid https URL",
			originalURL: "https://example.com",
			expectedErr: "",
		},

		{
			name:        "invalid URL",
			originalURL: "not-a-url",
			expectedErr: "original URL must be an HTTP or HTTPS URL",
		},
		{
			name:        "unsupported scheme",
			originalURL: "ftp://example.com",
			expectedErr: "original URL must be an HTTP or HTTPS URL",
		},
		{
			name:        "missing host",
			originalURL: "http:///path",
			expectedErr: "original URL must have a valid host",
		},
		{
			name:        "empty URL",
			originalURL: "",
			expectedErr: "original URL must be an HTTP or HTTPS URL",
		},
		{
			name:        "questionable URL",
			originalURL: "file:///tmp/x",
			expectedErr: "original URL must be an HTTP or HTTPS URL",
		},
	}

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			err := service.validateOriginalURL(subTest.originalURL)
			if subTest.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), subTest.expectedErr)
			}
		})
	}
}

func TestService_ValidateShortName(t *testing.T) {
	service := NewService(nil, nil)

	type subTest struct {
		name        string
		shortName   string
		expectedErr string
	}

	subTests := []subTest{
		{
			name:        "valid short name",
			shortName:   "abc123",
			expectedErr: "",
		},
		{
			name:        "too short",
			shortName:   "ab",
			expectedErr: "short name must be between 3 and 32 characters long",
		},
		{
			name:        "too long",
			shortName:   "a_very_long_short_name_exceeding_the_limit",
			expectedErr: "short name must be between 3 and 32 characters long",
		},
		{
			name:        "invalid characters",
			shortName:   "abc$123",
			expectedErr: "short name must match the pattern:",
		},
	}

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			err := service.validateShortName(subTest.shortName)
			if subTest.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), subTest.expectedErr)
			}
		})
	}
}

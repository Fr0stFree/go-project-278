package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToHashString(t *testing.T) {
	t.Run("should generate valid hash string", func(t *testing.T) {
		type testCase struct {
			name  string
			value string
			size  int
		}

		tests := []testCase{
			{
				name:  "should generate hash string",
				value: "https://example.com",
				size:  8,
			},
			{
				name:  "should generate single character",
				value: "https://example.com",
				size:  1,
			},
			{
				name:  "should handle empty value",
				value: "",
				size:  8,
			},
			{
				name:  "should handle zero size",
				value: "https://example.com",
				size:  0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := ToHashString(tt.value, tt.size)

				assert.Len(t, result, tt.size)

				for _, char := range result {
					assert.Contains(t, charset, string(char))
				}
			})
		}
	})

	t.Run("should generate consistent hash for same input", func(t *testing.T) {
		value := "https://example.com"
		size := 8

		hash1 := ToHashString(value, size)
		hash2 := ToHashString(value, size)

		assert.Equal(t, hash1, hash2)
	})

	t.Run("should generate different hashes for different inputs", func(t *testing.T) {
		size := 8

		hash1 := ToHashString("https://example.com/first", size)
		hash2 := ToHashString("https://example.com/second", size)

		assert.NotEqual(t, hash1, hash2)
	})
}

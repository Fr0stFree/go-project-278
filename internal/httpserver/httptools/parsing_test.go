package httptools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseQueryRange(t *testing.T) {
	t.Run("should accept JSON whitespace", func(t *testing.T) {
		from, to, err := ParseQueryRange("[0, 9]")

		require.NoError(t, err)
		assert.Equal(t, 0, from)
		assert.Equal(t, 9, to)
	})

	t.Run("should reject trailing data", func(t *testing.T) {
		_, _, err := ParseQueryRange("[0,9]anything")

		require.Error(t, err)
	})
}

func TestParseQuerySort(t *testing.T) {
	t.Run("should accept JSON whitespace", func(t *testing.T) {
		field, order, err := ParseQuerySort(`["id", "ASC"]`)

		require.NoError(t, err)
		assert.Equal(t, "id", field)
		assert.Equal(t, "ASC", order)
	})

	t.Run("should reject trailing data", func(t *testing.T) {
		_, _, err := ParseQuerySort(`["id","ASC"]anything`)

		require.Error(t, err)
	})
}

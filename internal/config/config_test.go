package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewParsesApplicationAndHTTPSettings(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/shortener")
	t.Setenv("HTTP_BASE_URL", "https://short.example.com")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "7s")

	cfg, err := New()

	require.NoError(t, err)
	assert.Equal(t, "https://short.example.com", cfg.HTTP.BaseURL)
	assert.Equal(t, 7*time.Second, cfg.App.ShutdownTimeout)
}

// Package config defines application, HTTP, and database settings.
package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Root groups all application configuration sections.
type Root struct {
	App      App
	HTTP     HTTP
	Database Database
}

// App contains settings used by business logic.
type App struct {
	BaseURL string `env:"APP_BASE_URL" envDefault:"http://localhost:8080"`
}

// HTTP contains server address and timeout settings.
type HTTP struct {
	Port            int           `env:"HTTP_PORT" envDefault:"8080"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout     time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"10s"`
	MaxBodySize     int64         `env:"HTTP_MAX_BODY_SIZE" envDefault:"16384"` // 16KiB
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"10s"`
	CORS            CORS
	Sentry          Sentry
}

// CORS contains Cross-Origin Resource Sharing settings.
type CORS struct {
	AllowOrigins []string      `env:"HTTP_CORS_ALLOW_ORIGINS" envSeparator:"," envDefault:"*"`
	MaxAge       time.Duration `env:"HTTP_CORS_MAX_AGE" envDefault:"12h"`
}

// Database contains PostgreSQL connection pool settings.
type Database struct {
	URL                   string        `env:"DATABASE_URL,required"`
	MaxOpenConnections    int           `env:"DB_MAX_OPEN_CONNECTIONS" envDefault:"10"`
	MaxIdleConnections    int           `env:"DB_MAX_IDLE_CONNECTIONS" envDefault:"5"`
	ConnectionMaxLifetime time.Duration `env:"DB_CONNECTION_MAX_LIFETIME" envDefault:"5m"`
}

// Sentry contains Sentry error tracking service settings.
type Sentry struct {
	IsEnabled    bool          `env:"SENTRY_ENABLED" envDefault:"false"`
	DSN          string        `env:"SENTRY_DSN"`
	Environment  string        `env:"SENTRY_ENVIRONMENT" envDefault:"development"`
	FlushTimeout time.Duration `env:"SENTRY_FLUSH_TIMEOUT" envDefault:"2s"`
}

// New returns the default local development configuration.
func New() (*Root, error) {
	// .env is optional. In production variables normally come directly from the environment.
	_ = godotenv.Load()

	config, err := env.ParseAs[Root]()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

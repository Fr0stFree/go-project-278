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
	DataBase DataBase
}

// App contains settings used by business logic.
type App struct {
	BaseURL string `env:"APP_BASE_URL"`
}

// HTTP contains server address and timeout settings.
type HTTP struct {
	Port         int           `env:"HTTP_PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
}

// DataBase contains PostgreSQL connection pool settings.
type DataBase struct {
	URL                   string        `env:"DATABASE_URL,required"`
	MaxOpenConnections    int           `env:"DB_MAX_OPEN_CONNECTIONS" envDefault:"10"`
	MaxIdleConnections    int           `env:"DB_MAX_IDLE_CONNECTIONS" envDefault:"5"`
	ConnectionMaxLifetime time.Duration `env:"DB_CONNECTION_MAX_LIFETIME" envDefault:"5m"`
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

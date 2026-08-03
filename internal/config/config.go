// Package config defines application, HTTP, and database settings.
package config

import "time"

// Root groups all application configuration sections.
type Root struct {
	App      App
	HTTP     HTTP
	DataBase DataBase
}

// App contains settings used by business logic.
type App struct {
	BaseURL string
}

// HTTP contains server address and timeout settings.
type HTTP struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DataBase contains PostgreSQL connection pool settings.
type DataBase struct {
	Host                  string
	Port                  int
	User                  string
	Password              string
	DBName                string
	IsSSLEnabled          bool
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
}

// New returns the default local development configuration.
func New() *Root {
	// TODO: make it configurable via environment variable or config file
	return &Root{
		HTTP: HTTP{
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		DataBase: DataBase{
			Host:                  "localhost",
			Port:                  5432,
			User:                  "postgres",
			Password:              "postgres",
			DBName:                "postgres",
			IsSSLEnabled:          false,
			MaxOpenConnections:    10,
			MaxIdleConnections:    5,
			ConnectionMaxLifetime: 5 * time.Minute,
		},
		App: App{
			BaseURL: "http://localhost:8080",
		},
	}
}

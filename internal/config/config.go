// Package config provides configuration structures for the application, including HTTP server settings.
package config

import "time"

// Root represents the overall configuration for the application.
type Root struct {
	App     App
	HTTP    HTTP
	Storage Storage
}

// App represents the configuration for business logic of the application.
type App struct {
	BaseURL string
}

// HTTP represents the configuration for the HTTP server.
type HTTP struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Storage represents the configuration for the storage backend.
type Storage struct {
	Type                  string
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

// New creates a new Config instance with default values.
func New() *Root {
	// TODO: make it configurable via environment variable or config file
	return &Root{
		HTTP: HTTP{
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		Storage: Storage{
			Type:               "postgres",
			Host:               "localhost",
			Port:               5432,
			User:               "postgres",
			Password:           "postgres",
			DBName:             "postgres",
			IsSSLEnabled:       false,
			MaxOpenConnections: 10,
			MaxIdleConnections: 5,
			ConnectionMaxLifetime:    5 * time.Minute,
		},
		App: App{
			BaseURL: "http://localhost:8080",
		},
	}
}

// Package config provides configuration structures for the application, including HTTP server settings.
package config

import "time"

// Config represents the overall configuration for the application.
type Config struct {
	HTTP HTTPConfig
}

// HTTPConfig represents the configuration for the HTTP server.
type HTTPConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// New creates a new Config instance with default values.
func New() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

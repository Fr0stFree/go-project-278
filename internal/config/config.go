package config

import "time"

type Config struct {
	HTTP HTTPConfig
}

type HTTPConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func New() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

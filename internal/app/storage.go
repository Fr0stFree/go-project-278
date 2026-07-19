package app

import (
	"shortener/internal/config"
	"shortener/internal/storage"
	"shortener/internal/storage/memory"
	"shortener/internal/storage/postgres"
)

func startLinkRepository(config *config.Storage) storage.AbstractLinkRepository {
	switch config.Type {
	case "memory":
		return memory.NewLinkRepository()
	case "postgres":
		repository, err := postgres.NewLinkRepository(config)
		if err != nil {
			panic("Failed to initialize PostgreSQL repository: " + err.Error())
		}
		return repository
	default:
		panic("Unsupported storage type: " + config.Type)
	}
}

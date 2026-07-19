package app

import (
	"shortener/internal/config"
	"shortener/internal/storage"
	"shortener/internal/storage/memory"
	"shortener/internal/storage/postgres"
)

func selectLinkRepository(config *config.Storage) storage.AbstractLinkRepository {
	switch config.Type {
	case "memory":
		return memory.NewLinkRepository()
	case "postgres":
		return postgres.NewLinkRepository(config)
	default:
		panic("Unsupported storage type: " + config.Type)
	}
}

package app
import (
	"shortener/internal/config"
	"shortener/internal/storage/memory"
	"shortener/internal/storage/postgres"
	"shortener/internal/storage"
)
func newLinkRepository(config config.StorageConfig) storage.AbstractLinkRepository {
	switch config.Type {
	case "memory":
		return memory.NewLinkRepository()
	case "postgres":
		return postgres.NewLinkRepository(config.DSN)
	default:
		panic("Unsupported storage type: " + config.Type)
	}
}
package app

import (
	"shortener/internal/config"
	"shortener/internal/database/postgres"
	"shortener/internal/database/repositories/link"
)

type combinedRepos struct {
	Link *link.Repository
}

func prepareRepos(cfg *config.DataBase) *combinedRepos {
	db, err := postgres.NewDataBase(cfg, &link.Record{})
	if err != nil {
		panic("Failed to initialize PostgreSQL repository: " + err.Error())
	}

	return &combinedRepos{Link: link.NewRepository(db)}
}

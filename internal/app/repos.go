package app

import (
	"shortener/internal/config"
	"shortener/internal/database/postgres"
	"shortener/internal/database/repositories/link"
	"shortener/internal/database/repositories/linkvisit"
)

type combinedRepos struct {
	link      *link.Repository
	linkVisit *linkvisit.Repository
}

func buildRepos(cfg *config.DataBase) *combinedRepos {
	db, err := postgres.NewDataBase(cfg, &link.Record{}, &linkvisit.Record{})
	if err != nil {
		panic("Failed to initialize PostgreSQL repository: " + err.Error())
	}

	return &combinedRepos{link: link.NewRepository(db), linkVisit: linkvisit.NewRepository(db)}
}

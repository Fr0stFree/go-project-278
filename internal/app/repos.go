package app

import (
	"fmt"

	"shortener/internal/config"
	"shortener/internal/database/postgres"
	"shortener/internal/database/repositories/link"
	"shortener/internal/database/repositories/linkvisit"
)

type combinedRepos struct {
	link      *link.Repository
	linkVisit *linkvisit.Repository
}

func buildRepos(cfg *config.DataBase) (*combinedRepos, error) {
	db, err := postgres.NewDataBase(cfg, &link.Record{}, &linkvisit.Record{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL repository: %w", err)
	}

	return &combinedRepos{link: link.NewRepository(db), linkVisit: linkvisit.NewRepository(db)}, nil
}

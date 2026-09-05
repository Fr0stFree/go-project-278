package app

import (
	"fmt"

	"shortener/internal/config"
	"shortener/internal/db"
	"shortener/internal/db/models/link"
	"shortener/internal/db/models/linkvisit"
)

type combinedRepos struct {
	link      *link.Repository
	linkVisit *linkvisit.Repository
}

func buildRepos(cfg *config.DataBase) (*combinedRepos, error) {
	database, err := db.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL repository: %w", err)
	}

	return &combinedRepos{link: link.NewRepository(database), linkVisit: linkvisit.NewRepository(database)}, nil
}

package app

import (
	"shortener/internal/config"
	"shortener/internal/services/shortener"
)

type combinedServices struct {
	shortener *shortener.Service
}

func buildServices(cfg *config.App, repos *combinedRepos) combinedServices {
	return combinedServices{
		shortener: shortener.NewService(repos.link, repos.linkVisit, cfg),
	}
}

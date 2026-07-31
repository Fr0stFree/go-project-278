package app

import (
	"shortener/internal/config"
	"shortener/internal/services/shortener"
)

type combinedServices struct {
	Shortener *shortener.Service
}

func prepareServices(cfg *config.App, repos *combinedRepos) combinedServices {
	return combinedServices{
		Shortener: shortener.NewService(repos.Link, cfg),
	}
}

package app

import (
	"shortener/internal/config"
	"shortener/internal/services/metrics"
	"shortener/internal/services/shortener"
)

type combinedServices struct {
	shortener *shortener.Service
	metrics   *metrics.Service
}

func prepareServices(cfg *config.App, repos *combinedRepos) combinedServices {
	return combinedServices{
		shortener: shortener.NewService(repos.link, cfg),
		metrics:   metrics.NewService(repos.linkVisit, cfg),
	}
}

package httpserver

import (
	"fmt"
	"net/http"
	"shortener/internal/config"
)

func New(cfg config.HTTPConfig) *http.Server {
	router := newRouter()
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}

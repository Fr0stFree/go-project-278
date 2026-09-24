package middleware

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"shortener/internal/config"
)

// CORS returns a middleware that handles Cross-Origin Resource Sharing requests.
func CORS(cfg config.CORS) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           cfg.MaxAge,
	})
}

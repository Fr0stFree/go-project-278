package httpserver

import (
	"github.com/gin-gonic/gin"
)

// NewRouter creates a Gin router with logging and recovery middleware.
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return router
}

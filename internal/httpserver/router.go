package httpserver

import (
	"github.com/gin-gonic/gin"
)

// NewRouter creates a new instance of the Gin router with default middleware applied.
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return router
}

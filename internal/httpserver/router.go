package httpserver

import (
	"github.com/gin-gonic/gin"
)

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/ping", pingHandler)

	return router
}

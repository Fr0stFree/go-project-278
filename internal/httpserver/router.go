package httpserver

import (
	"github.com/gin-gonic/gin"
)

func newRouter(c *combinedHandlers) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/ping", c.Health.Ping)
	router.POST("/api/links", c.Link.Create)
	router.GET("/api/links/:id", c.Link.Get)
	router.GET("/api/links", c.Link.List)
	router.PUT("/api/links/:id", c.Link.Update)
	router.DELETE("/api/links/:id", c.Link.Delete)

	return router
}

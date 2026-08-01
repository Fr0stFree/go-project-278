package health

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router gin.IRouter) {
	router.GET("/ping", ping)
}

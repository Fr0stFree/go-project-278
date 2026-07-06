package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	const Port int = 8080
	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})

	if err := router.Run(fmt.Sprintf(":%d", Port)); err != nil {
		panic(err)
	}
}
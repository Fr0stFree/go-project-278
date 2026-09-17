package httptools

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaxBodySize middleware limits the maximum size of the request body.
func MaxBodySize(maxBytes int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Body == nil {
			ctx.Next()
		}

		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxBytes)
		ctx.Next()
	}
}

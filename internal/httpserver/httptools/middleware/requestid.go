package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIDHeader = "X-Request-ID"
	requestIDCtxKey = "request_id"
)

// RequestID ensures that every request has a correlation identifier.
func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Set(requestIDCtxKey, requestID)
		ctx.Header(requestIDHeader, requestID)

		ctx.Next()
	}
}

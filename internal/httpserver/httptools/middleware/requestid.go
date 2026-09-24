package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"shortener/internal/httpserver/httptools/httpparam"
)

// RequestID ensures that every request has a correlation identifier.
func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := httpparam.ReadRequestIDHeader(ctx)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		httpparam.WriteRequestIDContext(ctx, requestID)
		httpparam.WriteRequestIDHeader(ctx, requestID)

		ctx.Next()
	}
}

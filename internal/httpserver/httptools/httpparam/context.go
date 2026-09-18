package httpparam

import "github.com/gin-gonic/gin"

const requestIDCtxKey = "request_id"

// ReadRequestIDContext retrieves the request ID from the gin context.
func ReadRequestIDContext(ctx *gin.Context) string {
	return ctx.GetString(requestIDCtxKey)
}

// ReadClientIP retrieves the client IP address from the gin context.
func ReadClientIP(ctx *gin.Context) string {
	return ctx.ClientIP()
}

// WriteRequestIDContext sets the request ID in the gin context.
func WriteRequestIDContext(ctx *gin.Context, requestID string) {
	ctx.Set(requestIDCtxKey, requestID)
}

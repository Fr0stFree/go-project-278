package middleware

import (
	"log/slog"
	"time"

	"shortener/internal/httpserver/httptools/httpparam"

	"github.com/gin-gonic/gin"
)

// AccessLogger logs HTTP requests and responses
func AccessLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		slog.InfoContext(
			ctx.Request.Context(),
			"HTTP request",
			slog.String("request_id", httpparam.ReadRequestIDContext(ctx)),
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.URL.Path),
			slog.Int("status", ctx.Writer.Status()),
			slog.String("client_ip", httpparam.ReadClientIP(ctx)),
			slog.Duration("duration", time.Since(start)),
			slog.Int64("content_length", ctx.Request.ContentLength),
			slog.Int("response_size", ctx.Writer.Size()),
			slog.String("user_agent", httpparam.ReadUserAgentHeader(ctx)),
		)
	}
}

// Package httpparam provides utilities for reading and writing HTTP parameters in the application.
package httpparam

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	requestIDHeader    = "X-Request-ID"
	contentRangeHeader = "Content-Range"
	userAgentHeader    = "User-Agent"
	referrerHeader     = "Referer"
)

// WriteContentRangeHeader formats a Content-Range header value using the actual number of returned records.
func WriteContentRangeHeader(ctx *gin.Context, unit string, from, count, totalCount int) {
	if count == 0 {
		ctx.Header(contentRangeHeader, fmt.Sprintf("%s */%d", unit, totalCount))

		return
	}

	to := from + count - 1
	ctx.Header(contentRangeHeader, fmt.Sprintf("%s %d-%d/%d", unit, from, to, totalCount))
}

// ReadRequestIDHeader retrieves the request ID from the request header.
func ReadRequestIDHeader(ctx *gin.Context) string {
	return ctx.GetHeader(requestIDHeader)
}

// WriteRequestIDHeader sets the request ID in the request header.
func WriteRequestIDHeader(ctx *gin.Context, requestID string) {
	ctx.Header(requestIDHeader, requestID)
}

// ReadUserAgentHeader retrieves the user agent from the request header.
func ReadUserAgentHeader(ctx *gin.Context) string {
	return ctx.GetHeader(userAgentHeader)
}

// ReadReferrerHeader retrieves the referrer from the request header.
func ReadReferrerHeader(ctx *gin.Context) string {
	return ctx.GetHeader(referrerHeader)
}

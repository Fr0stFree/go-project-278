// Package httpparam provides utilities for reading and writing HTTP parameters in the application.
package httpparam

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// WriteContentRangeHeader formats a Content-Range header value using the actual number of returned records.
func WriteContentRangeHeader(ctx *gin.Context, unit string, from, count, totalCount int) {
	if count == 0 {
		ctx.Header("Content-Range", fmt.Sprintf("%s */%d", unit, totalCount))

		return
	}

	to := from + count - 1
	ctx.Header("Content-Range", fmt.Sprintf("%s %d-%d/%d", unit, from, to, totalCount))
}

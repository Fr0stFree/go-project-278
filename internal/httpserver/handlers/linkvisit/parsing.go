package linkvisit

import (
	"fmt"
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

func parseFilterOpts(ctx *gin.Context) (*shortener.LinkVisitListOptionsBuilder, error) {
	builder := shortener.NewLinkVisitListOptionsBuilder()

	rangeRaw := ctx.Query("range")
	if rangeRaw != "" {
		var from, to int
		if _, err := fmt.Sscanf(rangeRaw, "[%d,%d]", &from, &to); err != nil {
			return nil, shortener.NewValidationError(fmt.Sprintf("invalid range format: %s", rangeRaw), "range")
		}

		builder.WithRange(from, to)
	}

	if builder.Error() != nil {
		return nil, builder.Error()
	}

	return builder, nil
}

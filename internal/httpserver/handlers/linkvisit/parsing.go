package linkvisit

import (
	"shortener/internal/httpserver"
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
)

func parseFilterOpts(ctx *gin.Context) (*shortener.LinkVisitListOptionsBuilder, error) {
	builder := shortener.NewLinkVisitListOptionsBuilder()

	rangeRaw := ctx.Query("range")
	if rangeRaw != "" {
		from, to, err := httpserver.ParseQueryRange(rangeRaw)
		if err != nil {
			return nil, shortener.NewValidationError(err.Error(), "range")
		}

		builder.WithRange(from, to)
	}

	return builder, nil
}

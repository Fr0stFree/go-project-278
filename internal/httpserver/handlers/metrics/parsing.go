package metrics

import (
	"fmt"
	"shortener/internal/database/repositories/linkvisit"

	"github.com/gin-gonic/gin"
)

func parseFilterOpts(ctx *gin.Context) (*linkvisit.FilterOpts, error) {
	opts := linkvisit.NewFilterOpts()

	rangeRaw := ctx.Query("range")
	if rangeRaw != "" {
		var from, to int
		if _, err := fmt.Sscanf(rangeRaw, "[%d,%d]", &from, &to); err != nil {
			return nil, err
		}

		if _, err := opts.WithRange(from, to); err != nil {
			return nil, err
		}
	}

	// TODO: Implement sorting parsing if needed in the future

	return opts, nil
}

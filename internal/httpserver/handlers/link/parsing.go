package link

import (
	"github.com/gin-gonic/gin"

	"shortener/internal/httpserver"
	"shortener/internal/services/shortener"
)

func parseLinkID(ctx *gin.Context) (uint, error) {
	linkID, err := httpserver.ParsePositiveIntParam(ctx.Param("id"))
	if err != nil {
		return 0, shortener.NewValidationError(err.Error(), "link_id")
	}

	return uint(linkID), nil
}

func parseFilterOpts(ctx *gin.Context) (*shortener.LinkListOptionsBuilder, error) {
	builder := shortener.NewLinkListOptionsBuilder()

	rangeRaw := ctx.Query("range")
	if rangeRaw != "" {
		from, to, err := httpserver.ParseQueryRange(rangeRaw)
		if err != nil {
			return nil, shortener.NewValidationError(err.Error(), "range")
		}

		builder.WithRange(from, to)
	}

	sortRaw := ctx.Query("sort")
	if sortRaw != "" {
		sortBy, sortOrder, err := httpserver.ParseQuerySort(sortRaw)
		if err != nil {
			return nil, shortener.NewValidationError(err.Error(), "sort")
		}

		builder.WithSort(sortBy, sortOrder)
	}

	return builder, nil
}

package link

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"shortener/internal/services/shortener"
)

func parseLinkID(ctx *gin.Context) (uint, error) {
	linkID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return 0, shortener.NewValidationError(fmt.Sprintf("invalid link ID: %s", ctx.Param("id")), "link_id")
	}

	if linkID < 0 {
		return 0, shortener.NewValidationError(fmt.Sprintf("invalid link ID: %d", linkID), "link_id")
	}

	return uint(linkID), nil
}

func parseFilterOpts(ctx *gin.Context) (*shortener.LinkListOptionsBuilder, error) {
	builder := shortener.NewLinkListOptionsBuilder()

	rangeRaw := ctx.Query("range")
	if rangeRaw != "" {
		var from, to int
		if _, err := fmt.Sscanf(rangeRaw, "[%d,%d]", &from, &to); err != nil {
			return nil, shortener.NewValidationError(fmt.Sprintf("invalid range format: %s", rangeRaw), "range")
		}

		builder.WithRange(from, to)
	}

	sortRaw := ctx.Query("sort")
	if sortRaw != "" {
		var sortBy, sortOrder string
		if _, err := fmt.Sscanf(sortRaw, "[%q,%q]", &sortBy, &sortOrder); err != nil {
			return nil, shortener.NewValidationError(fmt.Sprintf("invalid sort format: %s", sortRaw), "sort")
		}

		builder.WithSort(sortBy, sortOrder)
	}

	return builder, nil
}

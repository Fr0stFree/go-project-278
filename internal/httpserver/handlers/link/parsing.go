package link

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"shortener/internal/database/repositories/link"
)

func parseLinkID(ctx *gin.Context) (uint, error) {
	linkID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return 0, errors.New("invalid link ID")
	}

	if linkID < 0 {
		return 0, fmt.Errorf("link ID must be a non-negative integer, got %d", linkID)
	}

	return uint(linkID), nil
}

func parseFilterOpts(ctx *gin.Context) (*link.FilterOpts, error) {
	opts := link.NewFilterOpts()

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

	sortRaw := ctx.Query("sort")
	if sortRaw != "" {
		var sortBy, sortOrder string
		if _, err := fmt.Sscanf(sortRaw, "[%q,%q]", &sortBy, &sortOrder); err != nil {
			return nil, err
		}

		if _, err := opts.WithSort(sortBy, sortOrder); err != nil {
			return nil, err
		}
	}

	return opts, nil
}

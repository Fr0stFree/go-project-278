package link

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"shortener/internal/services/shortener"
)

func parseLinkID(ctx *gin.Context) (int, error) {
	linkID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return 0, errors.New("invalid link ID")
	}

	return linkID, nil
}

func parseLinksRange(ctx *gin.Context) (*shortener.LinkRange, error) {
	rangeRaw := ctx.Query("range")
	if rangeRaw == "" {
		return nil, nil
	}

	var LinkRange shortener.LinkRange

	_, err := fmt.Sscanf(rangeRaw, "[%d-%d]", &LinkRange.From, &LinkRange.To)
	if err != nil {
		return nil, errors.New("invalid range: expected format '[from-to]' with integers")
	}

	if !LinkRange.IsValid() {
		return nil, errors.New("invalid range: 'from' must be less than or equal to 'to', and both must be greater than or equal to 1")
	}

	return &LinkRange, nil
}

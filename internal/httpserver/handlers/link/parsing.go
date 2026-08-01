package link

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	"shortener/internal/services/shortener"
)

func parseLinksRange(ctx *gin.Context) (*shortener.LinkRange, error) {
	rangeRaw := ctx.Query("range")
	if rangeRaw == "" {
		return nil, nil
	}

	var LinkRange shortener.LinkRange

	_, err := fmt.Sscanf(rangeRaw, "[%d-%d]", &LinkRange.From, &LinkRange.To)
	if err != nil {
		return nil, err
	}

	if LinkRange.From < 0 || LinkRange.To < 0 || LinkRange.From > LinkRange.To {
		return nil, errors.New("invalid range values")
	}

	return &LinkRange, nil
}

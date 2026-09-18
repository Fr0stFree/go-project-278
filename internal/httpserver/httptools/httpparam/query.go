package httpparam

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// RangeQuery represents a range query parameter with a starting index and count.gj
type RangeQuery struct {
	From  int
	Count int
}

// ReadRangeQuery reads a range query parameter from the request context and returns a RangeQuery.
func ReadRangeQuery(ctx *gin.Context, paramName string) (*RangeQuery, error) {
	rangeRaw := ctx.Query(paramName)
	if rangeRaw == "" {
		return nil, nil
	}

	var values []int
	if err := json.Unmarshal([]byte(rangeRaw), &values); err != nil || len(values) != 2 {
		return nil, fmt.Errorf("invalid range format: %s", rangeRaw)
	}

	return &RangeQuery{
		From:  values[0],
		Count: values[1],
	}, nil
}

// SortQuery represents a sorting query parameter with a field and direction.
type SortQuery struct {
	Field     string
	Direction string
}

// ReadSortQuery reads a sort query parameter from the request context and returns a SortQuery .
func ReadSortQuery(ctx *gin.Context, paramName string) (*SortQuery, error) {
	sortRaw := ctx.Query(paramName)
	if sortRaw == "" {
		return nil, nil
	}

	var values []string
	if err := json.Unmarshal([]byte(sortRaw), &values); err != nil || len(values) != 2 {
		return nil, fmt.Errorf("invalid sort format: %s", sortRaw)
	}

	return &SortQuery{
		Field:     values[0],
		Direction: values[1],
	}, nil
}

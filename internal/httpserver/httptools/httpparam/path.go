package httpparam

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ReadStringPath reads a string from the named path parameter.
func ReadStringPath(ctx *gin.Context, paramName string) (string, error) {
	paramValue := ctx.Param(paramName)
	if paramValue == "" {
		return "", fmt.Errorf("missing path parameter: %s", paramName)
	}

	return paramValue, nil
}

// ReadNonNegativeIntPath reads a non-negative integer from the named path parameter.
func ReadNonNegativeIntPath(ctx *gin.Context, paramName string) (uint, error) {
	paramValue := ctx.Param(paramName)
	if paramValue == "" {
		return 0, fmt.Errorf("missing path parameter: %s", paramName)
	}

	value, err := strconv.Atoi(paramValue)
	if err != nil {
		return 0, fmt.Errorf("invalid positive integer: %s", paramValue)
	}

	if value < 0 {
		return 0, fmt.Errorf("value must be non-negative: %d", value)
	}

	return uint(value), nil
}

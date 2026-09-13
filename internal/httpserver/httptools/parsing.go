// Package httptools provides utility functions for parsing and processing HTTP.
package httptools

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ParseQueryRange parses a range query parameter in the format "[from,to]".
func ParseQueryRange(rangeRaw string) (int, int, error) {
	var values []int
	if err := json.Unmarshal([]byte(rangeRaw), &values); err != nil || len(values) != 2 {
		return 0, 0, fmt.Errorf("invalid range format: %s", rangeRaw)
	}

	return values[0], values[1], nil
}

// ParseQuerySort parses a sort query parameter in the format "[field,order]".
func ParseQuerySort(sortRaw string) (string, string, error) {
	var values []string
	if err := json.Unmarshal([]byte(sortRaw), &values); err != nil || len(values) != 2 {
		return "", "", fmt.Errorf("invalid sort format: %s", sortRaw)
	}

	return values[0], values[1], nil
}

// ParsePositiveIntParam parses a string parameter as a positive integer.
func ParsePositiveIntParam(param string) (int, error) {
	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("invalid positive integer: %s", param)
	}

	if value < 0 {
		return 0, fmt.Errorf("value must be non-negative: %d", value)
	}

	return value, nil
}

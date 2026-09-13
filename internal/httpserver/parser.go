package httpserver

import (
	"fmt"
	"regexp"
	"strconv"
)

var rangeRegexp = regexp.MustCompile(`^\[(\d+),(\d+)\]$`)

// ParseQueryRange parses a range query parameter in the format "[from,to]".
func ParseQueryRange(rangeRaw string) (int, int, error) {
	var (
		from, to int
		err      error
	)

	matches := rangeRegexp.FindStringSubmatch(rangeRaw)
	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("invalid range format: %s", rangeRaw)
	}

	from, err = strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid range start: %s", matches[1])
	}

	to, err = strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid range end: %s", matches[2])
	}

	return from, to, err
}

// ParseQuerySort parses a sort query parameter in the format "[field,order]".
func ParseQuerySort(sortRaw string) (string, string, error) {
	var sortBy, sortOrder string
	if _, err := fmt.Sscanf(sortRaw, "[%q,%q]", &sortBy, &sortOrder); err != nil {
		return "", "", fmt.Errorf("invalid sort format: %s", sortRaw)
	}

	return sortBy, sortOrder, nil
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

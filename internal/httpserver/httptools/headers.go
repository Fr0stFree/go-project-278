package httptools

import "fmt"

// FormatContentRange formats a Content-Range header value using the actual number of returned records.
func FormatContentRange(unit string, from, count, totalCount int) string {
	if count == 0 {
		return fmt.Sprintf("%s */%d", unit, totalCount)
	}

	to := from + count - 1

	return fmt.Sprintf("%s %d-%d/%d", unit, from, to, totalCount)
}

// Package utils provides utility functions for common operations.
package utils

import (
	"unicode"
)

// ToSnakeCase converts a string from camel case to snake case.
func ToSnakeCase(s string) string {
	var result []rune

	runes := []rune(s)

	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(runes[i-1]) ||
				(i+1 < len(runes) && unicode.IsLower(runes[i+1]))) {
				result = append(result, '_')
			}

			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}

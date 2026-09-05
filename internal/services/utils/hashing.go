package utils

import (
	"crypto/sha256"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// ToHashString generates a hash string of the given value with the specified size.
func ToHashString(value string, size int) string {
	hash := sha256.Sum256([]byte(value))

	shortCode := make([]byte, size)
	for i := range shortCode {
		index := int(hash[i]) % len(charset)
		shortCode[i] = charset[index]
	}

	return string(shortCode)
}

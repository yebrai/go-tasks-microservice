package id

import (
	"regexp"
	"strings"
)

// UUID regex pattern for validation
var uuidPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// IsValidUID validates if a string is a valid UUID format
func IsValidUID(value string) bool {
	if value == "" {
		return false
	}

	normalizedValue := strings.ToLower(value)

	return uuidPattern.MatchString(normalizedValue)
}

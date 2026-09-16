package utils

import (
	"regexp"
	"strings"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9]+`)
)

func MakeSlug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	// Replace non-alphanumeric with hyphen
	s = nonAlphanumericRegex.ReplaceAllString(s, "-")
	// Trim hyphens from ends
	s = strings.Trim(s, "-")
	return s
}

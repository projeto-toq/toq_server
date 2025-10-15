package utils

import "strings"

const wildcardToken = "*"

// NormalizeSearchPattern converts REST-friendly wildcards into SQL LIKE patterns.
// It trims spaces, replaces '*' with '%', and wraps plain terms with '%' if no wildcard is provided.
func NormalizeSearchPattern(input string) string {
	pattern := strings.TrimSpace(input)
	if pattern == "" {
		return ""
	}
	if strings.Contains(pattern, wildcardToken) {
		pattern = strings.ReplaceAll(pattern, wildcardToken, "%")
	}
	if !strings.ContainsAny(pattern, "%_") {
		pattern = "%" + pattern + "%"
	}
	return pattern
}

package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ParseCompletionForecast parses and validates completion forecast dates
//
// Accepts multiple formats and normalizes to YYYY-MM-DD for MySQL DATE column:
//   - "YYYY-MM-DD" (ISO 8601 date) → accepted as-is
//   - "YYYY-MM" (year-month) → normalized to "YYYY-MM-01" (first day of month)
//   - "YYYY-MM-DDT..." (RFC3339/ISO8601 timestamp) → extracts date portion
//
// Parameters:
//   - value: Raw string from client request (may contain timestamp, timezone, etc.)
//
// Returns:
//   - string: Normalized date in "YYYY-MM-DD" format ready for MySQL DATE column
//   - error: Validation error if format is invalid
//
// Example:
//
//	ParseCompletionForecast("2026-01-20T00:00:00Z") → "2026-01-20", nil
//	ParseCompletionForecast("2026-06") → "2026-06-01", nil
//	ParseCompletionForecast("invalid") → "", error
func ParseCompletionForecast(value string) (string, error) {
	// Trim whitespace
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("completion forecast cannot be empty")
	}

	var parsedDate time.Time
	var parseErr error

	// Try parsing as full RFC3339 timestamp (e.g., "2026-01-20T00:00:00Z")
	parsedDate, parseErr = time.Parse(time.RFC3339, trimmed)
	if parseErr == nil {
		// Extract date portion only (discard time and timezone)
		return parsedDate.Format("2006-01-02"), nil
	}

	// Try parsing as ISO8601 date with timezone offset (e.g., "2026-01-20+00:00")
	parsedDate, parseErr = time.Parse("2006-01-02T15:04:05Z07:00", trimmed)
	if parseErr == nil {
		return parsedDate.Format("2006-01-02"), nil
	}

	// Try parsing as date only (e.g., "2026-01-20")
	parsedDate, parseErr = time.Parse("2006-01-02", trimmed)
	if parseErr == nil {
		return parsedDate.Format("2006-01-02"), nil
	}

	// Try parsing as year-month only (e.g., "2026-06")
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}$`, trimmed)
	if matched {
		// Normalize to first day of month for MySQL DATE storage
		parsedDate, parseErr = time.Parse("2006-01", trimmed)
		if parseErr == nil {
			return parsedDate.Format("2006-01-02"), nil
		}
	}

	// No format matched
	return "", fmt.Errorf("invalid completion forecast format: expected YYYY-MM-DD, YYYY-MM, or RFC3339 timestamp")
}

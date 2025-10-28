package utils

import (
	"strings"
	"time"
)

const defaultTimezone = "America/Sao_Paulo"

// ResolveLocation loads a timezone location using an IANA identifier.
// When the provided value is blank the default application timezone is used.
func ResolveLocation(field, value string) (*time.Location, *HTTPError) {
	tz := strings.TrimSpace(value)
	if tz == "" {
		tz = defaultTimezone
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, ValidationError(field, "invalid timezone identifier")
	}

	return loc, nil
}

// NormalizeDateToLocationMidnight aligns the given timestamp to midnight on the provided location.
func NormalizeDateToLocationMidnight(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}

// NormalizeRangeToUTC converts optional range pointers to UTC using the provided location.
func NormalizeRangeToUTC(from, to *time.Time, loc *time.Location) (*time.Time, *time.Time) {
	var normalizedFrom *time.Time
	if from != nil && !from.IsZero() {
		f := NormalizeDateToLocationMidnight(*from, loc).UTC()
		normalizedFrom = &f
	}

	var normalizedTo *time.Time
	if to != nil && !to.IsZero() {
		t := NormalizeDateToLocationMidnight(*to, loc).UTC()
		normalizedTo = &t
	}

	return normalizedFrom, normalizedTo
}

// ConvertToUTC converts a timestamp expressed in any location to UTC keeping the instant.
func ConvertToUTC(value time.Time) time.Time {
	if value.IsZero() {
		return value
	}
	return value.In(time.UTC)
}

// ConvertToLocation converts a UTC or local timestamp to the provided location.
func ConvertToLocation(value time.Time, loc *time.Location) time.Time {
	if value.IsZero() {
		return value
	}
	return value.In(loc)
}

// FormatRFC3339InLocation returns the timestamp formatted with the given location applied.
func FormatRFC3339InLocation(value time.Time, loc *time.Location) string {
	if value.IsZero() {
		return ""
	}
	return ConvertToLocation(value, loc).Format(time.RFC3339)
}

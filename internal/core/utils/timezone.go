package utils

import (
	"strings"
	"time"
)

const defaultTimezone = "America/Sao_Paulo"

// ParseRFC3339Relaxed parses RFC3339 timestamps tolerating spaces in the offset portion.
// When the plus sign from the offset is converted into a space (e.g. by query decoding), this helper restores the signal.
func ParseRFC3339Relaxed(field, raw string) (time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return time.Time{}, ValidationError(field, field+" is required and must follow RFC3339")
	}

	normalized := normalizeRFC3339Offset(trimmed)
	value, err := time.Parse(time.RFC3339, normalized)
	if err != nil {
		return time.Time{}, ValidationError(field, field+" must be a valid RFC3339 timestamp")
	}
	return value, nil
}

func normalizeRFC3339Offset(value string) string {
	lastSpace := strings.LastIndex(value, " ")
	if lastSpace == -1 {
		return value
	}

	suffix := value[lastSpace+1:]
	if !looksLikeOffsetSuffix(suffix) {
		return value
	}

	prefix := value[:lastSpace]
	if len(suffix) == 0 {
		return value
	}

	if suffix[0] == '+' || suffix[0] == '-' {
		return prefix + suffix
	}

	return prefix + "+" + suffix
}

func looksLikeOffsetSuffix(s string) bool {
	if len(s) != 5 {
		return false
	}
	if s[2] != ':' {
		return false
	}
	first := s[0]
	if first == '+' || first == '-' {
		return true
	}
	return first >= '0' && first <= '9'
}

// DetermineRangeLocation selects the most appropriate location based on the provided timestamps.
// Priority: non-UTC "from", non-UTC "to", fallback location, otherwise UTC.
func DetermineRangeLocation(from, to time.Time, fallback *time.Location) *time.Location {
	if !from.IsZero() && from.Location() != time.UTC {
		return from.Location()
	}
	if !to.IsZero() && to.Location() != time.UTC {
		return to.Location()
	}
	if fallback != nil {
		return fallback
	}
	return time.UTC
}

// NormalizeRangeToUTC converts the provided timestamps into UTC using the supplied location.
// Zero values remain untouched.
func NormalizeRangeToUTC(from, to time.Time, loc *time.Location) (time.Time, time.Time) {
	if loc == nil {
		loc = time.UTC
	}
	if !from.IsZero() {
		from = ConvertToUTC(from.In(loc))
	}
	if !to.IsZero() {
		to = ConvertToUTC(to.In(loc))
	}
	return from, to
}

// NormalizePointerRangeToUTC converts pointer based timestamps into UTC versions using the location provided.
func NormalizePointerRangeToUTC(from, to *time.Time, loc *time.Location) (*time.Time, *time.Time) {
	if from == nil && to == nil {
		return nil, nil
	}
	var fromOut, toOut *time.Time
	if from != nil {
		convertedFrom, _ := NormalizeRangeToUTC(*from, time.Time{}, loc)
		fromOut = &convertedFrom
	}
	if to != nil {
		_, convertedTo := NormalizeRangeToUTC(time.Time{}, *to, loc)
		toOut = &convertedTo
	}
	return fromOut, toOut
}

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

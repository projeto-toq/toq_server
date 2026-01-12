package listingservices

import "time"

// monthsSince computes whole months between start and end timestamps using UTC calendar math.
func monthsSince(start, end time.Time) int {
	if start.IsZero() || end.IsZero() {
		return 0
	}
	if end.Before(start) {
		return 0
	}
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	total := years*12 + months
	if end.Day() < start.Day() {
		total--
	}
	if total < 0 {
		return 0
	}
	return total
}

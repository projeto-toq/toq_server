package scheduleservices

import (
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
)

const (
	availabilityDefaultLimit   = 100
	maxSlotDurationMinutes     = 12 * 60
	defaultSlotDurationMinutes = 60
)

var blockEntryTypes = map[schedulemodel.EntryType]struct{}{
	schedulemodel.EntryTypeBlock:          {},
	schedulemodel.EntryTypeTemporaryBlock: {},
}

type timeRange struct {
	start time.Time
	end   time.Time
}

func (r timeRange) isValid() bool {
	return r.start.Before(r.end)
}

func sanitizePagination(limit, page int) (int, int) {
	if limit <= 0 {
		limit = availabilityDefaultLimit
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func defaultSlotDuration(minutes uint16) time.Duration {
	if minutes == 0 || int(minutes) > maxSlotDurationMinutes {
		return time.Duration(defaultSlotDurationMinutes) * time.Minute
	}

	return time.Duration(minutes) * time.Minute
}

func subtractRange(base []timeRange, removal timeRange) []timeRange {
	if !removal.isValid() {
		return base
	}
	result := make([]timeRange, 0, len(base))
	for _, r := range base {
		if !r.isValid() {
			continue
		}
		if !removal.end.After(r.start) || !removal.start.Before(r.end) {
			result = append(result, r)
			continue
		}

		if removal.start.After(r.start) {
			left := timeRange{start: r.start, end: minTime(removal.start, r.end)}
			if left.isValid() {
				result = append(result, left)
			}
		}

		if removal.end.Before(r.end) {
			right := timeRange{start: maxTime(removal.end, r.start), end: r.end}
			if right.isValid() {
				result = append(result, right)
			}
		}
	}
	return result
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func splitIntoSlots(r timeRange, slot time.Duration) []AvailabilitySlot {
	slots := make([]AvailabilitySlot, 0)
	if !r.isValid() || slot <= 0 {
		return slots
	}

	cursor := r.start
	// Ranges are treated as half-open [start,end); allow slotEnd == r.end so the last slot ends exactly at the exclusive bound.
	for cursor.Add(slot).After(cursor) && !cursor.Add(slot).After(r.end) {
		slotEnd := cursor.Add(slot)
		slots = append(slots, AvailabilitySlot{StartsAt: cursor, EndsAt: slotEnd})
		cursor = slotEnd
	}

	return slots
}

func isBlockEntryType(entryType schedulemodel.EntryType) bool {
	_, ok := blockEntryTypes[entryType]
	return ok
}

func clampRange(r timeRange, min time.Time, max time.Time) (timeRange, bool) {
	start := maxTime(r.start, min)
	end := minTime(r.end, max)
	clamped := timeRange{start: start, end: end}
	if !clamped.isValid() {
		return timeRange{}, false
	}
	return clamped, true
}

func buildDayRange(day time.Time) timeRange {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	return timeRange{start: start, end: start.Add(24 * time.Hour)}
}

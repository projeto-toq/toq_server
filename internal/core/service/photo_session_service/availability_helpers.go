package photosessionservices

import (
	"sort"
	"time"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

const (
	slotIDPhotographerShift = 32
	slotIDTimeMask          = (1 << slotIDPhotographerShift) - 1
)

type timeRange struct {
	start time.Time
	end   time.Time
}

func encodeSlotID(photographerID uint64, start time.Time) uint64 {
	seconds := uint64(start.Unix()) & slotIDTimeMask
	return (photographerID << slotIDPhotographerShift) | seconds
}

func decodeSlotID(slotID uint64) (uint64, time.Time) {
	photographerID := slotID >> slotIDPhotographerShift
	seconds := int64(slotID & slotIDTimeMask)
	return photographerID, time.Unix(seconds, 0).UTC()
}

func clampRange(r timeRange, min, max time.Time) (timeRange, bool) {
	if r.end.Before(min) || r.start.After(max) {
		return timeRange{}, false
	}
	start := r.start
	if start.Before(min) {
		start = min
	}
	end := r.end
	if end.After(max) {
		end = max
	}
	if !end.After(start) {
		return timeRange{}, false
	}
	return timeRange{start: start, end: end}, true
}

func subtractRange(ranges []timeRange, removal timeRange) []timeRange {
	if !removal.end.After(removal.start) {
		return ranges
	}
	result := make([]timeRange, 0, len(ranges))
	for _, current := range ranges {
		if !current.end.After(removal.start) || !removal.end.After(current.start) {
			result = append(result, current)
			continue
		}
		if removal.start.After(current.start) {
			left := timeRange{start: current.start, end: removal.start}
			if left.end.After(left.start) {
				result = append(result, left)
			}
		}
		if removal.end.Before(current.end) {
			right := timeRange{start: removal.end, end: current.end}
			if right.end.After(right.start) {
				result = append(result, right)
			}
		}
	}
	return result
}

func splitIntoSlots(r timeRange, duration time.Duration) []timeRange {
	slots := make([]timeRange, 0)
	if duration <= 0 {
		return slots
	}
	start := r.start
	for start.Add(duration).Equal(r.end) || start.Add(duration).Before(r.end) {
		end := start.Add(duration)
		slots = append(slots, timeRange{start: start, end: end})
		start = end
	}
	return slots
}

func determineSlotPeriod(start time.Time) photosessionmodel.SlotPeriod {
	hour := start.Hour()
	if hour < 12 {
		return photosessionmodel.SlotPeriodMorning
	}
	return photosessionmodel.SlotPeriodAfternoon
}

func buildWorkingRanges(from, to time.Time, loc *time.Location, startHour, endHour int) []timeRange {
	ranges := make([]timeRange, 0)
	if endHour <= startHour {
		return ranges
	}

	startOfDay := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, loc)
	for day := startOfDay; day.Before(to); day = day.Add(24 * time.Hour) {
		start := time.Date(day.Year(), day.Month(), day.Day(), startHour, 0, 0, 0, loc)
		end := time.Date(day.Year(), day.Month(), day.Day(), endHour, 0, 0, 0, loc)
		if clamped, ok := clampRange(timeRange{start: start, end: end}, from, to); ok {
			ranges = append(ranges, clamped)
		}
	}

	return ranges
}

func applyBlockingEntries(ranges []timeRange, entries []photosessionmodel.AgendaEntryInterface, loc *time.Location) []timeRange {
	result := ranges
	for _, entry := range entries {
		if !entry.Blocking() {
			continue
		}
		removal := timeRange{
			start: entry.StartsAt().In(loc),
			end:   entry.EndsAt().In(loc),
		}
		if removal.end.After(removal.start) {
			result = subtractRange(result, removal)
		}
	}
	return result
}

func prunePastRanges(ranges []timeRange, reference time.Time) []timeRange {
	pruned := make([]timeRange, 0, len(ranges))
	for _, r := range ranges {
		if r.end.Before(reference) {
			continue
		}
		start := r.start
		if start.Before(reference) {
			start = reference
		}
		if r.end.After(start) {
			pruned = append(pruned, timeRange{start: start, end: r.end})
		}
	}
	return pruned
}

func sortAvailabilitySlots(slots []AvailabilitySlot, sortKey string) {
	switch sortKey {
	case "start_desc", "date_desc":
		sort.Slice(slots, func(i, j int) bool {
			if slots[i].Start.Equal(slots[j].Start) {
				return slots[i].PhotographerID > slots[j].PhotographerID
			}
			return slots[i].Start.After(slots[j].Start)
		})
	case "photographer_asc":
		sort.Slice(slots, func(i, j int) bool {
			if slots[i].PhotographerID == slots[j].PhotographerID {
				return slots[i].Start.Before(slots[j].Start)
			}
			return slots[i].PhotographerID < slots[j].PhotographerID
		})
	case "photographer_desc":
		sort.Slice(slots, func(i, j int) bool {
			if slots[i].PhotographerID == slots[j].PhotographerID {
				return slots[i].Start.Before(slots[j].Start)
			}
			return slots[i].PhotographerID > slots[j].PhotographerID
		})
	case "start_asc", "date_asc", "":
		sort.Slice(slots, func(i, j int) bool {
			if slots[i].Start.Equal(slots[j].Start) {
				return slots[i].PhotographerID < slots[j].PhotographerID
			}
			return slots[i].Start.Before(slots[j].Start)
		})
	default:
		sort.Slice(slots, func(i, j int) bool {
			if slots[i].Start.Equal(slots[j].Start) {
				return slots[i].PhotographerID < slots[j].PhotographerID
			}
			return slots[i].Start.Before(slots[j].Start)
		})
	}
}

func deriveSlotDuration(requested, configured int) time.Duration {
	base := configured
	if base <= 0 {
		base = 120
	}
	if requested > 0 {
		if requested != base {
			// Enforce single duration to keep slotID encoding (photographer+start) consistent.
			return time.Duration(-1)
		}
		base = requested
	}
	return time.Duration(base) * time.Minute
}

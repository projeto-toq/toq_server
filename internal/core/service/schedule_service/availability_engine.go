package scheduleservices

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// dailyAvailability groups free ranges for a day anchored on agenda timezone.
type dailyAvailability struct {
	dayStart time.Time
	ranges   []timeRange
}

// availabilityEngineCompute builds paginated free slots using agenda timezone and rules/entries.
func availabilityEngineCompute(requestRange schedulemodel.ScheduleRange, loc *time.Location, slotMinutes uint16, pagination schedulemodel.PaginationConfig, data schedulemodel.AvailabilityData) AvailabilityResult {
	fromLocal := utils.ConvertToLocation(requestRange.From, loc)
	toLocal := utils.ConvertToLocation(requestRange.To, loc)

	convertEntriesToLocation(data.Entries, loc)

	days := buildInitialAvailability(fromLocal, toLocal)
	applyRules(days, data.Rules, fromLocal, toLocal)
	applyEntries(days, data.Entries, fromLocal, toLocal)

	slotDuration := defaultSlotDuration(slotMinutes)
	allSlots := collectSlots(days, fromLocal, toLocal, slotDuration)
	total := len(allSlots)

	limit, offset := sanitizePagination(pagination.Limit, pagination.Page)
	if offset >= total {
		return AvailabilityResult{Slots: []AvailabilitySlot{}, Total: total, Timezone: loc.String()}
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return AvailabilityResult{Slots: allSlots[offset:end], Total: total, Timezone: loc.String()}
}

// availabilityEngineFits checks if a requested slot is fully contained in current free ranges.
func availabilityEngineFits(slot schedulemodel.ScheduleRange, loc *time.Location, data schedulemodel.AvailabilityData) bool {
	fromLocal := utils.ConvertToLocation(slot.From, loc)
	toLocal := utils.ConvertToLocation(slot.To, loc)

	convertEntriesToLocation(data.Entries, loc)

	days := buildInitialAvailability(fromLocal, toLocal)
	applyRules(days, data.Rules, fromLocal, toLocal)
	applyEntries(days, data.Entries, fromLocal, toLocal)

	return windowFits(days, timeRange{start: fromLocal, end: toLocal})
}

// buildAvailabilityRepoFilter normalizes the requested range to UTC for repository access.
func buildAvailabilityRepoFilter(filter schedulemodel.AvailabilityFilter, slot schedulemodel.ScheduleRange, loc *time.Location) schedulemodel.AvailabilityFilter {
	repoFilter := filter
	repoFilter.Range.From, repoFilter.Range.To = utils.NormalizeRangeToUTC(slot.From, slot.To, loc)
	repoFilter.Range.Loc = time.UTC
	return repoFilter
}

// resolveAgendaLocation loads the agenda timezone, defaulting to the global timezone when empty.
func resolveAgendaLocation(agenda schedulemodel.AgendaInterface) (*time.Location, *utils.HTTPError) {
	return utils.ResolveLocation("timezone", agenda.Timezone())
}

// mapAgendaError normalizes agenda lookup errors into domain-aware responses.
func mapAgendaError(ctx context.Context, logger *slog.Logger, err error, listingIdentityID int64) error {
	if errors.Is(err, sql.ErrNoRows) {
		return utils.NotFoundError("Agenda")
	}
	utils.SetSpanError(ctx, err)
	if logger != nil {
		logger.Error("schedule.agenda_error", "listing_identity_id", listingIdentityID, "err", err)
	}
	return utils.InternalError("")
}

// buildInitialAvailability seeds the daily availability ranges for the requested interval.
func buildInitialAvailability(from, to time.Time) []*dailyAvailability {
	startOfDay := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	days := make([]*dailyAvailability, 0)

	for day := startOfDay; day.Before(to); day = day.Add(24 * time.Hour) {
		dayRange := buildDayRange(day)
		if clamped, ok := clampRange(dayRange, from, to); ok {
			days = append(days, &dailyAvailability{
				dayStart: dayRange.start,
				ranges:   []timeRange{clamped},
			})
		}
	}

	return days
}

// applyRules removes blocked ranges configured by agenda rules.
// Rules must already be normalized to half-open semantics ([start,end)), meaning end_minute in DB is inclusive
// and was converted to an exclusive upper bound before this step.
func applyRules(days []*dailyAvailability, rules []schedulemodel.AgendaRuleInterface, from, to time.Time) {
	for _, rule := range rules {
		if !rule.IsActive() || rule.RuleType() != schedulemodel.RuleTypeBlock {
			continue
		}
		for _, day := range days {
			if day.dayStart.Weekday() != rule.DayOfWeek() {
				continue
			}
			block := timeRange{
				start: day.dayStart.Add(time.Duration(rule.StartMinutes()) * time.Minute),
				end:   day.dayStart.Add(time.Duration(rule.EndMinutes()) * time.Minute),
			}
			if clamped, ok := clampRange(block, day.dayStart, day.dayStart.Add(24*time.Hour)); ok {
				if final, ok := clampRange(clamped, from, to); ok {
					day.ranges = subtractRange(day.ranges, final)
				}
			}
		}
	}
}

// applyEntries removes blocked ranges produced by agenda entries.
func applyEntries(days []*dailyAvailability, entries []schedulemodel.AgendaEntryInterface, from, to time.Time) {
	dayLength := 24 * time.Hour
	for _, entry := range entries {
		if !entry.Blocking() {
			continue
		}

		removal := timeRange{start: entry.StartsAt(), end: entry.EndsAt()}
		for _, day := range days {
			dayRange := timeRange{start: day.dayStart, end: day.dayStart.Add(dayLength)}
			if clipped, ok := clampRange(removal, dayRange.start, dayRange.end); ok {
				if final, ok := clampRange(clipped, from, to); ok {
					day.ranges = subtractRange(day.ranges, final)
				}
			}
		}
	}
}

// collectSlots splits remaining ranges into slots respecting the requested duration.
func collectSlots(days []*dailyAvailability, from, to time.Time, slotDuration time.Duration) []AvailabilitySlot {
	slots := make([]AvailabilitySlot, 0)
	for _, day := range days {
		for _, free := range day.ranges {
			if clamped, ok := clampRange(free, from, to); ok {
				subSlots := splitIntoSlots(clamped, slotDuration)
				slots = append(slots, subSlots...)
			}
		}
	}
	return slots
}

// windowFits checks whether the target range is fully contained in available ranges.
func windowFits(days []*dailyAvailability, target timeRange) bool {
	for _, day := range days {
		for _, free := range day.ranges {
			if target.start.Equal(free.start) || target.start.After(free.start) {
				if target.end.Equal(free.end) || target.end.Before(free.end) {
					return true
				}
			}
		}
	}
	return false
}

// convertEntriesToLocation updates entry timestamps to the provided timezone.
func convertEntriesToLocation(entries []schedulemodel.AgendaEntryInterface, loc *time.Location) {
	for _, entry := range entries {
		entry.SetStartsAt(utils.ConvertToLocation(entry.StartsAt(), loc))
		entry.SetEndsAt(utils.ConvertToLocation(entry.EndsAt(), loc))
	}
}

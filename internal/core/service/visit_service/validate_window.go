package visitservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// validateWindow enforces lead time, horizon, and availability before creating a visit.
func (s *visitService) validateWindow(ctx context.Context, tx *sql.Tx, agenda schedulemodel.AgendaInterface, input CreateVisitInput) error {
	if err := s.validateLeadTimeAndHorizon(input); err != nil {
		return err
	}

	return s.validateAvailability(ctx, tx, agenda, input)
}

func (s *visitService) validateLeadTimeAndHorizon(input CreateVisitInput) error {
	now := time.Now().UTC()
	leadLimit := now.Add(time.Duration(s.config.MinHoursAhead) * time.Hour)
	if !input.ScheduledStart.After(leadLimit) {
		return utils.ValidationError("scheduledStart", fmt.Sprintf("must be at least %d hours in advance", s.config.MinHoursAhead))
	}

	horizonLimit := now.Add(time.Duration(s.config.MaxDaysAhead) * 24 * time.Hour)
	if input.ScheduledStart.After(horizonLimit) || input.ScheduledEnd.After(horizonLimit) {
		return utils.ValidationError("scheduledStart", fmt.Sprintf("must be within %d days from now", s.config.MaxDaysAhead))
	}

	return nil
}

func (s *visitService) validateAvailability(ctx context.Context, tx *sql.Tx, agenda schedulemodel.AgendaInterface, input CreateVisitInput) error {
	logger := utils.LoggerFromContext(ctx)
	loc, tzErr := utils.ResolveLocation("timezone", agenda.Timezone())
	if tzErr != nil {
		return tzErr
	}

	duration := input.ScheduledEnd.Sub(input.ScheduledStart)
	if duration <= 0 {
		return utils.ValidationError("scheduledTime", "end must be after start")
	}

	fromLocal := utils.ConvertToLocation(input.ScheduledStart, loc)
	toLocal := utils.ConvertToLocation(input.ScheduledEnd, loc)

	filter := schedulemodel.AvailabilityFilter{
		ListingIdentityID:  input.ListingIdentityID,
		Range:              schedulemodel.ScheduleRange{From: fromLocal, To: toLocal, Loc: time.UTC},
		SlotDurationMinute: uint16(duration.Minutes()),
	}
	filter.Range.From, filter.Range.To = utils.NormalizeRangeToUTC(fromLocal, toLocal, loc)

	data, err := s.scheduleRepo.GetAvailabilityData(ctx, tx, filter)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("visit.validate.availability.repo_error", "listing_identity_id", input.ListingIdentityID, "err", err)
		return utils.InternalError("")
	}

	for _, entry := range data.Entries {
		entry.SetStartsAt(utils.ConvertToLocation(entry.StartsAt(), loc))
		entry.SetEndsAt(utils.ConvertToLocation(entry.EndsAt(), loc))
	}

	days := buildInitialAvailability(fromLocal, toLocal)
	applyRules(days, data.Rules, fromLocal, toLocal)
	applyEntries(days, data.Entries, fromLocal, toLocal)

	if !windowFits(days, timeRange{start: fromLocal, end: toLocal}) {
		return utils.ConflictError("Requested slot is not available")
	}

	return nil
}

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

// Helpers below mirror the availability calculations from schedule_service.

type dailyAvailability struct {
	dayStart time.Time
	ranges   []timeRange
}

type timeRange struct {
	start time.Time
	end   time.Time
}

func (r timeRange) isValid() bool {
	return r.start.Before(r.end)
}

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

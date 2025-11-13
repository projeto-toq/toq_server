package scheduleservices

import (
	"context"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

type dailyAvailability struct {
	dayStart time.Time
	ranges   []timeRange
}

func (s *scheduleService) GetAvailability(ctx context.Context, filter schedulemodel.AvailabilityFilter) (AvailabilityResult, error) {
	if filter.ListingIdentityID <= 0 {
		return AvailabilityResult{}, utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}
	if filter.Range.From.IsZero() || filter.Range.To.IsZero() {
		return AvailabilityResult{}, utils.ValidationError("range", "from and to must be provided")
	}
	if err := validateRange(filter.Range.From, filter.Range.To); err != nil {
		return AvailabilityResult{}, err
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return AvailabilityResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.get_availability.tx_start_error", "err", txErr)
		return AvailabilityResult{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("schedule.get_availability.tx_rollback_error", "err", rbErr)
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, filter.ListingIdentityID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.get_availability.get_agenda_error", "listing_identity_id", filter.ListingIdentityID, "err", err)
		return AvailabilityResult{}, utils.InternalError("")
	}

	agendaLoc, tzErr := utils.ResolveLocation("timezone", agenda.Timezone())
	if tzErr != nil {
		return AvailabilityResult{}, tzErr
	}
	loc := filter.Range.Loc
	if loc == nil {
		loc = agendaLoc
	}

	repoFilter := filter
	repoFilter.Range.From, repoFilter.Range.To = utils.NormalizeRangeToUTC(filter.Range.From, filter.Range.To, loc)
	repoFilter.Range.Loc = time.UTC

	data, err := s.scheduleRepo.GetAvailabilityData(ctx, tx, repoFilter)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.get_availability.repo_error", "listing_identity_id", filter.ListingIdentityID, "err", err)
		return AvailabilityResult{}, utils.InternalError("")
	}

	fromLocal := utils.ConvertToLocation(filter.Range.From, loc)
	toLocal := utils.ConvertToLocation(filter.Range.To, loc)

	for _, entry := range data.Entries {
		entry.SetStartsAt(utils.ConvertToLocation(entry.StartsAt(), loc))
		entry.SetEndsAt(utils.ConvertToLocation(entry.EndsAt(), loc))
	}

	days := buildInitialAvailability(fromLocal, toLocal)
	applyRules(days, data.Rules, fromLocal, toLocal)
	applyEntries(days, data.Entries, fromLocal, toLocal)

	slotDuration := defaultSlotDuration(filter.SlotDurationMinute)
	allSlots := collectSlots(days, fromLocal, toLocal, slotDuration)
	total := len(allSlots)

	limit, offset := sanitizePagination(filter.Pagination.Limit, filter.Pagination.Page)
	if offset >= total {
		return AvailabilityResult{Slots: []AvailabilitySlot{}, Total: total, Timezone: loc.String()}, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return AvailabilityResult{Slots: allSlots[offset:end], Total: total, Timezone: loc.String()}, nil
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

package scheduleservices

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) ListAgendaEntries(ctx context.Context, filter schedulemodel.AgendaDetailFilter) (schedulemodel.AgendaDetailResult, error) {
	if filter.OwnerID <= 0 {
		return schedulemodel.AgendaDetailResult{}, utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}
	if filter.ListingID <= 0 {
		return schedulemodel.AgendaDetailResult{}, utils.ValidationError("listingId", "listingId must be greater than zero")
	}
	if err := validateRange(filter.Range.From, filter.Range.To); err != nil {
		return schedulemodel.AgendaDetailResult{}, err
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return schedulemodel.AgendaDetailResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.list_agenda_entries.tx_start_error", "err", txErr, "listing_id", filter.ListingID)
		return schedulemodel.AgendaDetailResult{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("schedule.list_agenda_entries.tx_rollback_error", "err", rbErr, "listing_id", filter.ListingID)
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingID(ctx, tx, filter.ListingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schedulemodel.AgendaDetailResult{}, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.list_agenda_entries.get_agenda_error", "listing_id", filter.ListingID, "err", err)
		return schedulemodel.AgendaDetailResult{}, utils.InternalError("")
	}

	if agenda.OwnerID() != filter.OwnerID {
		return schedulemodel.AgendaDetailResult{}, utils.AuthorizationError("Owner does not match listing agenda")
	}

	agendaLoc, tzErr := utils.ResolveLocation("timezone", agenda.Timezone())
	if tzErr != nil {
		return schedulemodel.AgendaDetailResult{}, tzErr
	}
	loc := filter.Range.Loc
	if loc == nil {
		loc = agendaLoc
	}

	fromUTC, toUTC := utils.NormalizeRangeToUTC(filter.Range.From, filter.Range.To, loc)
	localRange := schedulemodel.ScheduleRange{From: utils.ConvertToLocation(filter.Range.From, loc), To: utils.ConvertToLocation(filter.Range.To, loc)}
	fromLocal, toLocal := resolveTimelineRange(localRange, loc)

	entries, err := s.scheduleRepo.ListEntriesBetween(ctx, tx, agenda.ID(), fromUTC, toUTC)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.list_agenda_entries.list_entries_error", "listing_id", filter.ListingID, "err", err)
		return schedulemodel.AgendaDetailResult{}, utils.InternalError("")
	}

	for _, entry := range entries {
		entry.SetStartsAt(utils.ConvertToLocation(entry.StartsAt(), loc))
		entry.SetEndsAt(utils.ConvertToLocation(entry.EndsAt(), loc))
	}

	rules, err := s.scheduleRepo.ListRulesByAgenda(ctx, tx, agenda.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.list_agenda_entries.rules_error", "listing_id", filter.ListingID, "err", err)
		return schedulemodel.AgendaDetailResult{}, utils.InternalError("")
	}

	timeline := composeTimeline(entries, rules, fromLocal, toLocal, loc)
	sort.SliceStable(timeline, func(i, j int) bool {
		if timeline[i].StartsAt.Equal(timeline[j].StartsAt) {
			return timeline[i].EndsAt.Before(timeline[j].EndsAt)
		}
		return timeline[i].StartsAt.Before(timeline[j].StartsAt)
	})

	limit, offset := sanitizePagination(filter.Pagination.Limit, filter.Pagination.Page)
	total := len(timeline)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	items := make([]schedulemodel.AgendaTimelineItem, 0, end-offset)
	items = append(items, timeline[offset:end]...)

	return schedulemodel.AgendaDetailResult{Items: items, Total: int64(total), Timezone: loc.String()}, nil
}

func resolveTimelineRange(rng schedulemodel.ScheduleRange, loc *time.Location) (time.Time, time.Time) {
	const defaultTimelineDays = 14
	var from, to time.Time
	if rng.From.IsZero() {
		from = time.Now().In(loc)
	} else {
		from = rng.From.In(loc)
	}
	if rng.To.IsZero() {
		to = from.Add(defaultTimelineDays * 24 * time.Hour)
	} else {
		to = rng.To.In(loc)
	}
	if !from.Before(to) {
		to = from.Add(24 * time.Hour)
	}
	return from, to
}

func composeTimeline(entries []schedulemodel.AgendaEntryInterface, rules []schedulemodel.AgendaRuleInterface, from, to time.Time, loc *time.Location) []schedulemodel.AgendaTimelineItem {
	items := make([]schedulemodel.AgendaTimelineItem, 0)
	dayCursor := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, loc)
	endBoundary := to

	for _, entry := range entries {
		items = append(items, schedulemodel.AgendaTimelineItem{
			Source:    schedulemodel.TimelineSourceEntry,
			Entry:     entry,
			StartsAt:  entry.StartsAt(),
			EndsAt:    entry.EndsAt(),
			Weekday:   entry.StartsAt().In(loc).Weekday(),
			Recurring: false,
			Blocking:  entry.Blocking(),
		})
	}

	for day := dayCursor; day.Before(endBoundary); day = day.Add(24 * time.Hour) {
		for _, rule := range rules {
			if rule == nil || !rule.IsActive() {
				continue
			}
			if rule.RuleType() != schedulemodel.RuleTypeBlock {
				continue
			}
			if day.Weekday() != rule.DayOfWeek() {
				continue
			}
			window := buildRuleWindow(day, RuleTimeRange{StartMinute: rule.StartMinutes(), EndMinute: rule.EndMinutes()}, loc)
			if clamped, ok := clampRange(window, from, to); ok {
				items = append(items, schedulemodel.AgendaTimelineItem{
					Source:    schedulemodel.TimelineSourceRule,
					Rule:      rule,
					StartsAt:  clamped.start,
					EndsAt:    clamped.end,
					Weekday:   day.Weekday(),
					Recurring: true,
					Blocking:  true,
				})
			}
		}
	}

	return items
}

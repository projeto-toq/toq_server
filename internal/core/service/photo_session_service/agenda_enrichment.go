package photosessionservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

type photographerLocation struct {
	city  string
	state string
}

func (s *photoSessionService) loadPhotographerLocation(ctx context.Context, tx *sql.Tx, photographerID uint64) (photographerLocation, error) {
	user, err := s.userRepo.GetUserByID(ctx, tx, int64(photographerID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return photographerLocation{}, utils.NotFoundError("Photographer")
		}
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("service.list_agenda.photographer_get_error", "photographer_id", photographerID, "err", err)
		return photographerLocation{}, derrors.Wrap(err, derrors.KindInfra, "failed to load photographer")
	}

	loc := photographerLocation{
		city:  strings.TrimSpace(user.GetCity()),
		state: strings.ToUpper(strings.TrimSpace(user.GetState())),
	}
	return loc, nil
}

func (s *photoSessionService) fetchHolidaySlots(
	ctx context.Context,
	photographerID uint64,
	loc *time.Location,
	profile photographerLocation,
	from time.Time,
	to time.Time,
	occupied map[string]struct{},
) ([]AgendaSlot, error) {
	calendars, err := s.listCalendarsByLocation(ctx, profile)
	if err != nil {
		return nil, err
	}
	if len(calendars) == 0 {
		return nil, nil
	}

	logger := utils.LoggerFromContext(ctx)
	fromLocal := utils.ConvertToLocation(from, loc)
	toLocal := utils.ConvertToLocation(to, loc)
	if !toLocal.After(fromLocal) {
		toLocal = fromLocal.Add(time.Minute)
	}
	fromUTC := fromLocal.UTC()
	toUTC := toLocal.UTC()

	slots := make([]AgendaSlot, 0)
	for _, calendar := range calendars {
		filter := holidaymodel.CalendarDatesFilter{
			CalendarID: calendar.ID(),
			From:       &fromUTC,
			To:         &toUTC,
			Timezone:   loc.String(),
			Limit:      200,
			Page:       1,
		}

		for {
			result, err := s.holidayService.ListCalendarDates(ctx, filter)
			if err != nil {
				utils.SetSpanError(ctx, err)
				logger.Error("service.list_agenda.holiday_repo_error", "calendar_id", calendar.ID(), "err", err)
				return nil, derrors.Wrap(err, derrors.KindInfra, "failed to list holiday dates")
			}

			for _, date := range result.Dates {
				day := utils.ConvertToLocation(date.HolidayDate(), loc)
				start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, loc)
				end := start.Add(24 * time.Hour)

				if clamped, ok := clampRange(timeRange{start: start, end: end}, fromLocal, toLocal); ok {
					key := agendaSlotKey(photosessionmodel.AgendaEntryTypeHoliday, photosessionmodel.AgendaEntrySourceHoliday, clamped.start, clamped.end)
					if _, exists := occupied[key]; exists {
						continue
					}

					label := strings.TrimSpace(date.Label())
					if label == "" {
						label = "Holiday"
					}
					slot := newSyntheticHolidaySlot(photographerID, calendar.ID(), []string{label}, clamped.start, clamped.end, loc)
					occupied[key] = struct{}{}
					slots = append(slots, slot)
				}
			}

			if len(result.Dates) < filter.Limit {
				break
			}
			filter.Page++
		}
	}

	return slots, nil
}

func (s *photoSessionService) listCalendarsByLocation(ctx context.Context, profile photographerLocation) ([]holidaymodel.CalendarInterface, error) {
	active := true
	limit := 100
	logger := utils.LoggerFromContext(ctx)

	type scopeRequest struct {
		scope holidaymodel.CalendarScope
		state *string
		city  *string
	}

	scopes := []scopeRequest{{scope: holidaymodel.ScopeNational}}
	if profile.state != "" {
		state := profile.state
		scopes = append(scopes, scopeRequest{scope: holidaymodel.ScopeState, state: &state})
		if profile.city != "" {
			city := profile.city
			scopes = append(scopes, scopeRequest{scope: holidaymodel.ScopeCity, state: &state, city: &city})
		}
	}

	calendarMap := make(map[uint64]holidaymodel.CalendarInterface)
	for _, req := range scopes {
		page := 1
		for {
			filter := holidaymodel.CalendarListFilter{
				Scope:      &req.scope,
				OnlyActive: &active,
				Page:       page,
				Limit:      limit,
			}
			if req.state != nil {
				filter.State = req.state
			}
			if req.city != nil {
				filter.City = req.city
			}

			result, err := s.holidayService.ListCalendars(ctx, filter)
			if err != nil {
				utils.SetSpanError(ctx, err)
				logger.Error("service.list_agenda.calendar_list_error", "scope", req.scope, "err", err)
				return nil, derrors.Wrap(err, derrors.KindInfra, "failed to list holiday calendars")
			}

			for _, calendar := range result.Calendars {
				calendarMap[calendar.ID()] = calendar
			}

			if len(result.Calendars) < limit || len(result.Calendars) == 0 {
				break
			}
			page++
		}
	}

	calendars := make([]holidaymodel.CalendarInterface, 0, len(calendarMap))
	for _, calendar := range calendarMap {
		calendars = append(calendars, calendar)
	}

	return calendars, nil
}

func (s *photoSessionService) buildNonWorkingSlots(
	photographerID uint64,
	loc *time.Location,
	from time.Time,
	to time.Time,
	occupied map[string]struct{},
) []AgendaSlot {
	fromLocal := utils.ConvertToLocation(from, loc)
	toLocal := utils.ConvertToLocation(to, loc)
	if !toLocal.After(fromLocal) {
		toLocal = fromLocal.Add(time.Minute)
	}

	slots := make([]AgendaSlot, 0)
	dayCursor := time.Date(fromLocal.Year(), fromLocal.Month(), fromLocal.Day(), 0, 0, 0, 0, loc)
	limitDay := time.Date(toLocal.Year(), toLocal.Month(), toLocal.Day(), 0, 0, 0, 0, loc)
	if !toLocal.Equal(limitDay) {
		limitDay = limitDay.Add(24 * time.Hour)
	}

	for day := dayCursor; day.Before(limitDay); day = day.Add(24 * time.Hour) {
		baseRange := timeRange{start: day, end: day.Add(24 * time.Hour)}
		clampedDay, ok := clampRange(baseRange, fromLocal, toLocal)
		if !ok {
			continue
		}

		if day.Weekday() == time.Saturday || day.Weekday() == time.Sunday {
			key := agendaSlotKey(photosessionmodel.AgendaEntryTypeBlock, photosessionmodel.AgendaEntrySourceOnboarding, clampedDay.start, clampedDay.end)
			if _, exists := occupied[key]; !exists {
				slot := newSyntheticBlockSlot(photographerID, clampedDay.start, clampedDay.end, loc, "Weekend")
				occupied[key] = struct{}{}
				slots = append(slots, slot)
			}
			continue
		}

		businessStart := time.Date(day.Year(), day.Month(), day.Day(), s.cfg.BusinessStartHour, 0, 0, 0, loc)
		businessEnd := time.Date(day.Year(), day.Month(), day.Day(), s.cfg.BusinessEndHour, 0, 0, 0, loc)

		if businessStart.After(day) {
			if earlyRange, ok := clampRange(timeRange{start: day, end: businessStart}, fromLocal, toLocal); ok && earlyRange.end.After(earlyRange.start) {
				key := agendaSlotKey(photosessionmodel.AgendaEntryTypeBlock, photosessionmodel.AgendaEntrySourceOnboarding, earlyRange.start, earlyRange.end)
				if _, exists := occupied[key]; !exists {
					slot := newSyntheticBlockSlot(photographerID, earlyRange.start, earlyRange.end, loc, "Outside business hours")
					occupied[key] = struct{}{}
					slots = append(slots, slot)
				}
			}
		}

		if businessEnd.Before(day.Add(24 * time.Hour)) {
			if lateRange, ok := clampRange(timeRange{start: businessEnd, end: day.Add(24 * time.Hour)}, fromLocal, toLocal); ok && lateRange.end.After(lateRange.start) {
				key := agendaSlotKey(photosessionmodel.AgendaEntryTypeBlock, photosessionmodel.AgendaEntrySourceOnboarding, lateRange.start, lateRange.end)
				if _, exists := occupied[key]; !exists {
					slot := newSyntheticBlockSlot(photographerID, lateRange.start, lateRange.end, loc, "Outside business hours")
					occupied[key] = struct{}{}
					slots = append(slots, slot)
				}
			}
		}
	}

	return slots
}

func syntheticSlotID(namespace string, start time.Time, extras ...string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(namespace))
	h.Write([]byte(start.UTC().Format(time.RFC3339Nano)))
	for _, extra := range extras {
		h.Write([]byte{0})
		h.Write([]byte(extra))
	}
	return h.Sum64() | (uint64(1) << 63)
}

func agendaSlotKey(entryType photosessionmodel.AgendaEntryType, source photosessionmodel.AgendaEntrySource, start, end time.Time) string {
	builder := strings.Builder{}
	builder.WriteString(string(entryType))
	builder.WriteByte('|')
	builder.WriteString(string(source))
	builder.WriteByte('|')
	builder.WriteString(strconv.FormatInt(start.UTC().UnixNano(), 10))
	builder.WriteByte('|')
	builder.WriteString(strconv.FormatInt(end.UTC().UnixNano(), 10))
	return builder.String()
}

func newSyntheticBlockSlot(photographerID uint64, start, end time.Time, loc *time.Location, reason string) AgendaSlot {
	trimmed := strings.TrimSpace(reason)
	if trimmed == "" {
		trimmed = "Outside business hours"
	}
	normalized := strings.ReplaceAll(strings.ToLower(trimmed), " ", "-")
	groupID := fmt.Sprintf("synthetic-%s-%s", normalized, start.Format("2006-01-02"))

	return AgendaSlot{
		EntryID:        syntheticSlotID("block", start, strconv.FormatUint(photographerID, 10), normalized),
		PhotographerID: photographerID,
		EntryType:      photosessionmodel.AgendaEntryTypeBlock,
		Source:         photosessionmodel.AgendaEntrySourceOnboarding,
		Start:          start,
		End:            end,
		Status:         string(photosessionmodel.SlotStatusBlocked),
		GroupID:        groupID,
		Reason:         trimmed,
		Timezone:       loc.String(),
	}
}

func newSyntheticHolidaySlot(photographerID uint64, calendarID uint64, labels []string, start, end time.Time, loc *time.Location) AgendaSlot {
	cleanLabels := make([]string, 0, len(labels))
	for _, label := range labels {
		trimmed := strings.TrimSpace(label)
		if trimmed != "" {
			cleanLabels = append(cleanLabels, trimmed)
		}
	}
	if len(cleanLabels) == 0 {
		cleanLabels = append(cleanLabels, "Holiday")
	}

	slot := AgendaSlot{
		EntryID:            syntheticSlotID("holiday", start, strconv.FormatUint(calendarID, 10)),
		PhotographerID:     photographerID,
		EntryType:          photosessionmodel.AgendaEntryTypeHoliday,
		Source:             photosessionmodel.AgendaEntrySourceHoliday,
		SourceID:           calendarID,
		Start:              start,
		End:                end,
		Status:             string(photosessionmodel.SlotStatusBlocked),
		GroupID:            fmt.Sprintf("holiday-%s", start.Format("2006-01-02")),
		IsHoliday:          true,
		HolidayLabels:      cleanLabels,
		HolidayCalendarIDs: []uint64{calendarID},
		Timezone:           loc.String(),
		Reason:             cleanLabels[0],
	}

	return slot
}

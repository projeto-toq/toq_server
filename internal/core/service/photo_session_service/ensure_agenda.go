package photosessionservices

import (
	"context"
	"database/sql"
	"sort"
	"strings"
	"time"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// EnsurePhotographerAgenda provisions bootstrap agenda entries for a photographer.
func (s *photoSessionService) EnsurePhotographerAgenda(ctx context.Context, input EnsureAgendaInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.ensure_agenda.tx_start_error", "err", err)
		return utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
				utils.SetSpanError(ctx, rollbackErr)
				logger.Error("photo_session.ensure_agenda.tx_rollback_error", "err", rollbackErr)
			}
		}
	}()

	if err := s.ensurePhotographerAgendaInternal(ctx, tx, input); err != nil {
		return err
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.ensure_agenda.tx_commit_error", "err", err)
		return utils.InternalError("")
	}
	committed = true
	return nil
}

// EnsurePhotographerAgendaWithTx provisions bootstrap agenda entries using an existing transaction.
func (s *photoSessionService) EnsurePhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error {
	if tx == nil {
		return utils.InternalError("")
	}
	ctx = utils.ContextWithLogger(ctx)
	return s.ensurePhotographerAgendaInternal(ctx, tx, input)
}

// RefreshPhotographerAgenda re-applies bootstrap agenda rules for the photographer.
func (s *photoSessionService) RefreshPhotographerAgenda(ctx context.Context, input EnsureAgendaInput) error {
	return s.EnsurePhotographerAgenda(ctx, input)
}

// RefreshPhotographerAgendaWithTx re-applies bootstrap agenda rules using an existing transaction.
func (s *photoSessionService) RefreshPhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error {
	return s.EnsurePhotographerAgendaWithTx(ctx, tx, input)
}

func (s *photoSessionService) ensurePhotographerAgendaInternal(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error {
	logger := utils.LoggerFromContext(ctx)

	if input.PhotographerID == 0 {
		return utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}

	loc, tzErr := resolveLocation(input.Timezone)
	if tzErr != nil {
		return tzErr
	}

	workdayStart := input.WorkdayStartHour
	if workdayStart <= 0 {
		workdayStart = defaultWorkdayStartHour
	}
	workdayEnd := input.WorkdayEndHour
	if workdayEnd <= 0 {
		workdayEnd = defaultWorkdayEndHour
	}
	if workdayEnd <= workdayStart {
		return utils.ValidationError("workdayEndHour", "workdayEndHour must be greater than workdayStartHour")
	}

	horizonMonths := input.HorizonMonths
	if horizonMonths <= 0 {
		horizonMonths = defaultHorizonMonths
	}

	windowStart := s.now().In(loc).Truncate(24 * time.Hour)
	windowEnd := windowStart.AddDate(0, horizonMonths, 0)

	if s.holidayRepo != nil && input.HolidayCalendarID != nil {
		ids := make([]uint64, 0, 1)
		if *input.HolidayCalendarID > 0 {
			ids = append(ids, *input.HolidayCalendarID)
		}
		if err := s.holidayRepo.ReplaceAssociations(ctx, tx, input.PhotographerID, ids); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("photo_session.ensure_agenda.replace_holiday_assoc_error", "photographer_id", input.PhotographerID, "err", err)
			return utils.InternalError("")
		}
	}

	if _, err := s.repo.DeleteEntriesBySource(ctx, tx, input.PhotographerID, photosessionmodel.AgendaEntryTypeBlock, photosessionmodel.AgendaEntrySourceOnboarding, nil); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.ensure_agenda.delete_block_entries_error", "photographer_id", input.PhotographerID, "err", err)
		return utils.InternalError("")
	}

	if _, err := s.repo.DeleteEntriesBySource(ctx, tx, input.PhotographerID, photosessionmodel.AgendaEntryTypeHoliday, photosessionmodel.AgendaEntrySourceHoliday, nil); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.ensure_agenda.delete_holiday_entries_error", "photographer_id", input.PhotographerID, "err", err)
		return utils.InternalError("")
	}

	bootstrapEntries := make([]photosessionmodel.AgendaEntryInterface, 0)
	bootstrapEntries = append(bootstrapEntries, s.buildDefaultBlockEntries(input.PhotographerID, loc, windowStart, windowEnd, workdayStart, workdayEnd)...)

	holidayEntries, err := s.buildHolidayEntries(ctx, input.PhotographerID, loc, windowStart, windowEnd, input.HolidayCalendarID)
	if err != nil {
		return err
	}
	bootstrapEntries = append(bootstrapEntries, holidayEntries...)

	if len(bootstrapEntries) == 0 {
		return nil
	}

	if _, err := s.repo.CreateEntries(ctx, tx, bootstrapEntries); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.ensure_agenda.create_entries_error", "photographer_id", input.PhotographerID, "err", err)
		return utils.InternalError("")
	}

	return nil
}

func (s *photoSessionService) buildDefaultBlockEntries(photographerID uint64, loc *time.Location, windowStart, windowEnd time.Time, workdayStart, workdayEnd int) []photosessionmodel.AgendaEntryInterface {
	entries := make([]photosessionmodel.AgendaEntryInterface, 0)
	for day := windowStart; day.Before(windowEnd); day = day.AddDate(0, 0, 1) {
		dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, loc)
		dayEnd := dayStart.Add(24 * time.Hour)

		if dayStart.Weekday() == time.Saturday || dayStart.Weekday() == time.Sunday {
			entries = append(entries, newBlockingEntry(photographerID, loc.String(), dayStart.UTC(), dayEnd.UTC(), "Weekend"))
			continue
		}

		if workdayStart > 0 {
			blockStart := dayStart
			blockEnd := time.Date(day.Year(), day.Month(), day.Day(), workdayStart, 0, 0, 0, loc)
			if blockEnd.After(blockStart) {
				entries = append(entries, newBlockingEntry(photographerID, loc.String(), blockStart.UTC(), blockEnd.UTC(), "Outside business hours"))
			}
		}

		if workdayEnd < 24 {
			blockStart := time.Date(day.Year(), day.Month(), day.Day(), workdayEnd, 0, 0, 0, loc)
			blockEnd := dayEnd
			if blockEnd.After(blockStart) {
				entries = append(entries, newBlockingEntry(photographerID, loc.String(), blockStart.UTC(), blockEnd.UTC(), "Outside business hours"))
			}
		}
	}
	return entries
}

func (s *photoSessionService) buildHolidayEntries(ctx context.Context, photographerID uint64, loc *time.Location, windowStart, windowEnd time.Time, calendarID *uint64) ([]photosessionmodel.AgendaEntryInterface, error) {
	if calendarID == nil || *calendarID == 0 {
		return nil, nil
	}

	dates, err := s.fetchHolidayDates(ctx, []uint64{*calendarID}, windowStart, windowEnd)
	if err != nil {
		return nil, err
	}

	entries := make([]photosessionmodel.AgendaEntryInterface, 0, len(dates))
	for _, item := range dates {
		day := item.date.In(loc)
		start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, loc)
		end := start.Add(24 * time.Hour)

		entry := photosessionmodel.NewAgendaEntry()
		entry.SetPhotographerUserID(photographerID)
		entry.SetEntryType(photosessionmodel.AgendaEntryTypeHoliday)
		entry.SetSource(photosessionmodel.AgendaEntrySourceHoliday)
		entry.SetSourceID(item.calendarID)
		entry.SetStartsAt(start.UTC())
		entry.SetEndsAt(end.UTC())
		entry.SetBlocking(true)
		entry.SetTimezone(loc.String())
		reason := strings.TrimSpace(item.label)
		if reason == "" {
			reason = "Holiday"
		}
		entry.SetReason(reason)
		entries = append(entries, entry)
	}

	return entries, nil
}

func newBlockingEntry(photographerID uint64, timezone string, start, end time.Time, reason string) photosessionmodel.AgendaEntryInterface {
	entry := photosessionmodel.NewAgendaEntry()
	entry.SetPhotographerUserID(photographerID)
	entry.SetEntryType(photosessionmodel.AgendaEntryTypeBlock)
	entry.SetSource(photosessionmodel.AgendaEntrySourceOnboarding)
	entry.SetStartsAt(start)
	entry.SetEndsAt(end)
	entry.SetBlocking(true)
	entry.SetTimezone(timezone)
	if trimmed := strings.TrimSpace(reason); trimmed != "" {
		entry.SetReason(trimmed)
	}
	return entry
}

type holidayDateInfo struct {
	calendarID uint64
	date       time.Time
	label      string
}

func (s *photoSessionService) fetchHolidayDates(ctx context.Context, calendarIDs []uint64, from, to time.Time) ([]holidayDateInfo, error) {
	if len(calendarIDs) == 0 {
		return nil, nil
	}

	fromUTC := from.UTC()
	toUTC := to.UTC()
	if toUTC.Before(fromUTC) {
		toUTC = fromUTC
	}

	entries := make([]holidayDateInfo, 0)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	for _, calendarID := range calendarIDs {
		if calendarID == 0 {
			continue
		}

		filter := holidaymodel.CalendarDatesFilter{
			CalendarID: calendarID,
			From:       &fromUTC,
			To:         &toUTC,
			Limit:      200,
			Page:       1,
		}

		for {
			result, err := s.holidayService.ListCalendarDates(ctx, filter)
			if err != nil {
				utils.SetSpanError(ctx, err)
				logger.Error("photo_session.ensure_agenda.holiday_list_error", "calendar_id", calendarID, "err", err)
				return nil, utils.InternalError("")
			}

			for _, date := range result.Dates {
				entries = append(entries, holidayDateInfo{
					calendarID: calendarID,
					date:       date.HolidayDate(),
					label:      date.Label(),
				})
			}

			if len(result.Dates) < filter.Limit {
				break
			}
			filter.Page++
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].date.Equal(entries[j].date) {
			return entries[i].calendarID < entries[j].calendarID
		}
		return entries[i].date.Before(entries[j].date)
	})

	return entries, nil
}

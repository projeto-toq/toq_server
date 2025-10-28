package photosessionservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	photosessionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/photo_session_repository"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	holidayservices "github.com/projeto-toq/toq_server/internal/core/service/holiday_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	defaultHorizonMonths    = 3
	defaultWorkdayStartHour = 8
	defaultWorkdayEndHour   = 19
	defaultTimezone         = "America/Sao_Paulo"
	maxTimeOffReasonLength  = 255
	defaultAgendaPage       = 1
	defaultAgendaSize       = 20
	maxAgendaPageSize       = 100
)

// AgendaEntrySource identifies the origin of an agenda entry.
type AgendaEntrySource string

const (
	AgendaEntrySourceSlot    AgendaEntrySource = "SLOT"
	AgendaEntrySourceHoliday AgendaEntrySource = "HOLIDAY"
	AgendaEntrySourceTimeOff AgendaEntrySource = "TIME_OFF"
)

// PhotoSessionServiceInterface exposes orchestration helpers around photographer slots.
type PhotoSessionServiceInterface interface {
	EnsurePhotographerAgenda(ctx context.Context, input EnsureAgendaInput) error
	EnsurePhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error
	RefreshPhotographerAgenda(ctx context.Context, input EnsureAgendaInput) error
	RefreshPhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error
	CreateTimeOff(ctx context.Context, input TimeOffInput) (uint64, error)
	CreateTimeOffWithTx(ctx context.Context, tx *sql.Tx, input TimeOffInput) (uint64, error)
	DeleteTimeOff(ctx context.Context, input DeleteTimeOffInput) error
	DeleteTimeOffWithTx(ctx context.Context, tx *sql.Tx, input DeleteTimeOffInput) error
	ListAgenda(ctx context.Context, input ListAgendaInput) (ListAgendaOutput, error)
	UpdateSessionStatus(ctx context.Context, input UpdateSessionStatusInput) error
}

type photoSessionService struct {
	repo           photosessionrepository.PhotoSessionRepositoryInterface
	listingRepo    listingrepository.ListingRepoPortInterface
	holidayService holidayservices.HolidayServiceInterface
	globalService  globalservice.GlobalServiceInterface
	cfg            Config
	now            func() time.Time
}

// NewPhotoSessionService wires a photo session service with its dependencies.
func NewPhotoSessionService(
	repo photosessionrepository.PhotoSessionRepositoryInterface,
	listingRepo listingrepository.ListingRepoPortInterface,
	holidayService holidayservices.HolidayServiceInterface,
	globalService globalservice.GlobalServiceInterface,
) PhotoSessionServiceInterface {
	return NewPhotoSessionServiceWithConfig(repo, listingRepo, holidayService, globalService, Config{})
}

// NewPhotoSessionServiceWithConfig wires a photo session service with explicit config.
func NewPhotoSessionServiceWithConfig(
	repo photosessionrepository.PhotoSessionRepositoryInterface,
	listingRepo listingrepository.ListingRepoPortInterface,
	holidayService holidayservices.HolidayServiceInterface,
	globalService globalservice.GlobalServiceInterface,
	cfg Config,
) PhotoSessionServiceInterface {
	return &photoSessionService{
		repo:           repo,
		listingRepo:    listingRepo,
		holidayService: holidayService,
		globalService:  globalService,
		cfg:            normalizeConfig(cfg),
		now:            time.Now,
	}
}

// EnsurePhotographerAgenda provisions and normalizes future slots for a photographer.
func (s *photoSessionService) EnsurePhotographerAgenda(ctx context.Context, input EnsureAgendaInput) error {
	prepared, err := s.prepareEnsureContext(input)
	if err != nil {
		return err
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		return txErr
	}

	committed := false
	defer func() {
		if !committed {
			_ = s.globalService.RollbackTransaction(ctx, tx)
		}
	}()

	if err = s.ensurePhotographerAgendaWithPrepared(ctx, tx, input, prepared); err != nil {
		return err
	}

	if err = s.globalService.CommitTransaction(ctx, tx); err != nil {
		return err
	}
	committed = true

	return nil
}

// EnsurePhotographerAgendaWithTx provisions slots using an existing transaction.
func (s *photoSessionService) EnsurePhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error {
	prepared, err := s.prepareEnsureContext(input)
	if err != nil {
		return err
	}
	return s.ensurePhotographerAgendaWithPrepared(ctx, tx, input, prepared)
}

// RefreshPhotographerAgenda re-runs the ensure workflow to extend the horizon.
func (s *photoSessionService) RefreshPhotographerAgenda(ctx context.Context, input EnsureAgendaInput) error {
	return s.EnsurePhotographerAgenda(ctx, input)
}

// RefreshPhotographerAgendaWithTx re-runs ensure workflow with provided transaction.
func (s *photoSessionService) RefreshPhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error {
	return s.EnsurePhotographerAgendaWithTx(ctx, tx, input)
}

// CreateTimeOff registers a new time-off entry and re-syncs slots.
func (s *photoSessionService) CreateTimeOff(ctx context.Context, input TimeOffInput) (uint64, error) {
	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		return 0, txErr
	}

	committed := false
	defer func() {
		if !committed {
			_ = s.globalService.RollbackTransaction(ctx, tx)
		}
	}()

	id, err := s.CreateTimeOffWithTx(ctx, tx, input)
	if err != nil {
		return 0, err
	}

	if err = s.globalService.CommitTransaction(ctx, tx); err != nil {
		return 0, err
	}
	committed = true

	return id, nil
}

// CreateTimeOffWithTx registers a new time-off entry within an existing transaction.
func (s *photoSessionService) CreateTimeOffWithTx(ctx context.Context, tx *sql.Tx, input TimeOffInput) (uint64, error) {
	if err := validateTimeOffInput(input); err != nil {
		return 0, err
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	id, err := s.repo.CreateTimeOff(ctx, tx, input.toModel())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.time_off.create_error", "err", err)
		return 0, utils.InternalError("")
	}

	ensureInput := EnsureAgendaInput{
		PhotographerID:    input.PhotographerID,
		Timezone:          input.Timezone,
		HolidayCalendarID: input.HolidayCalendarID,
		HorizonMonths:     input.HorizonMonths,
		WorkdayStartHour:  input.WorkdayStartHour,
		WorkdayEndHour:    input.WorkdayEndHour,
	}
	if err = s.ensurePhotographerAgendaWithPrepared(ctx, tx, ensureInput, nil); err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteTimeOff removes an existing time-off entry and re-syncs slots.
func (s *photoSessionService) DeleteTimeOff(ctx context.Context, input DeleteTimeOffInput) error {
	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		return txErr
	}

	committed := false
	defer func() {
		if !committed {
			_ = s.globalService.RollbackTransaction(ctx, tx)
		}
	}()

	if err := s.DeleteTimeOffWithTx(ctx, tx, input); err != nil {
		return err
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		return err
	}
	committed = true

	return nil
}

// DeleteTimeOffWithTx removes a time-off entry within an existing transaction.
func (s *photoSessionService) DeleteTimeOffWithTx(ctx context.Context, tx *sql.Tx, input DeleteTimeOffInput) error {
	if input.TimeOffID == 0 {
		return utils.ValidationError("timeOffId", "timeOffId must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err = s.repo.DeleteTimeOff(ctx, tx, input.TimeOffID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Time off")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.time_off.delete_error", "err", err, "time_off_id", input.TimeOffID)
		return utils.InternalError("")
	}

	ensureInput := EnsureAgendaInput{
		PhotographerID:    input.PhotographerID,
		Timezone:          input.Timezone,
		HolidayCalendarID: input.HolidayCalendarID,
		HorizonMonths:     input.HorizonMonths,
		WorkdayStartHour:  input.WorkdayStartHour,
		WorkdayEndHour:    input.WorkdayEndHour,
	}
	return s.ensurePhotographerAgendaWithPrepared(ctx, tx, ensureInput, nil)
}

// UpdateSessionStatus updates the status of a photo session for a photographer.
func (s *photoSessionService) UpdateSessionStatus(ctx context.Context, input UpdateSessionStatusInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.SessionID == 0 {
		return derrors.Validation("sessionId must be greater than zero", map[string]any{"sessionId": "greater_than_zero"})
	}
	if input.PhotographerID == 0 {
		return derrors.Auth("unauthorized")
	}

	statusStr := strings.ToUpper(strings.TrimSpace(input.Status))
	if statusStr == "" {
		return derrors.Validation("status is required", map[string]any{"status": "required"})
	}

	if statusStr != string(photosessionmodel.BookingStatusAccepted) && statusStr != string(photosessionmodel.BookingStatusRejected) {
		return derrors.BadRequest("status must be ACCEPTED or REJECTED")
	}

	status := photosessionmodel.BookingStatus(statusStr)

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.tx_start_error", "err", err)
		return derrors.Infra("failed to start transaction", err)
	}

	committed := false
	defer func() {
		if !committed {
			if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
				utils.SetSpanError(ctx, rollbackErr)
				logger.Error("photo_session.update_status.tx_rollback_error", "err", rollbackErr)
			}
		}
	}()

	booking, err := s.repo.GetBookingForUpdate(ctx, tx, input.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("session not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.get_booking_error", "err", err, "session_id", input.SessionID)
		return derrors.Infra("failed to load session booking", err)
	}

	slot, err := s.repo.GetSlotForUpdate(ctx, tx, booking.SlotID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("slot not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.get_slot_error", "err", err, "slot_id", booking.SlotID())
		return derrors.Infra("failed to load session slot", err)
	}

	if slot.PhotographerUserID() != input.PhotographerID {
		return derrors.Forbidden("session does not belong to photographer")
	}

	if booking.Status() != photosessionmodel.BookingStatusPendingApproval {
		return derrors.Conflict("session is not pending approval")
	}

	if err = s.repo.UpdateBookingStatus(ctx, tx, booking.ID(), status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("session not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.update_error", "err", err, "session_id", booking.ID())
		return derrors.Infra("failed to update session status", err)
	}

	if err = s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.tx_commit_error", "err", err)
		return derrors.Infra("failed to commit transaction", err)
	}
	committed = true

	logger.Info("photo_session.status.updated", "session_id", booking.ID(), "slot_id", slot.ID(), "photographer_id", input.PhotographerID, "status", statusStr)

	return nil
}

// ListAgenda retrieves the photographer's agenda within a specified date range.
func (s *photoSessionService) ListAgenda(ctx context.Context, input ListAgendaInput) (ListAgendaOutput, error) {
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "service.ListAgenda")
	if err != nil {
		return ListAgendaOutput{}, derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	if err := validateListAgendaInput(input); err != nil {
		return ListAgendaOutput{}, err
	}

	page := input.Page
	if page <= 0 {
		page = defaultAgendaPage
	}

	size := input.Size
	if size <= 0 {
		size = defaultAgendaSize
	}
	if size > maxAgendaPageSize {
		size = maxAgendaPageSize
	}

	offset := (page - 1) * size

	// Use a transaction for a consistent read, although it's a read-only operation.
	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.LoggerFromContext(ctx).Error("service.list_agenda.tx_start_error", "err", txErr)
		utils.SetSpanError(ctx, txErr)
		return ListAgendaOutput{}, derrors.Wrap(txErr, derrors.KindInfra, "failed to start transaction")
	}
	defer func() {
		_ = s.globalService.RollbackTransaction(ctx, tx) // Always rollback a read-only transaction
	}()

	slots, total, err := s.repo.ListSlotsByRangePaginated(ctx, tx, input.PhotographerID, input.StartDate.UTC(), input.EndDate.UTC(), size, offset)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("service.list_agenda.repo_error", "err", err, "photographer_id", input.PhotographerID)
		return ListAgendaOutput{}, derrors.Wrap(err, derrors.KindInfra, "failed to list agenda slots")
	}

	loc, locErr := resolveLocation(input.Timezone)
	if locErr != nil {
		return ListAgendaOutput{}, locErr
	}

	holidayInfos, holidayErr := s.fetchHolidayDates(ctx, input.HolidayCalendarIDs, input.StartDate, input.EndDate)
	if holidayErr != nil {
		return ListAgendaOutput{}, holidayErr
	}
	holidayByDate := make(map[string][]holidayDateInfo, len(holidayInfos))
	for _, info := range holidayInfos {
		key := dateKeyFromTime(info.date.In(loc))
		holidayByDate[key] = append(holidayByDate[key], info)
	}
	matchedHolidayDates := make(map[string]struct{})
	resultTimezone := loc.String()

	timeOffRecords, timeOffErr := s.repo.ListTimeOff(ctx, tx, input.PhotographerID, input.StartDate.UTC(), input.EndDate.UTC())
	if timeOffErr != nil {
		utils.SetSpanError(ctx, timeOffErr)
		utils.LoggerFromContext(ctx).Error("service.list_agenda.time_off_error", "err", timeOffErr, "photographer_id", input.PhotographerID)
		return ListAgendaOutput{}, derrors.Wrap(timeOffErr, derrors.KindInfra, "failed to list agenda time off")
	}

	entries := make([]AgendaSlot, 0, len(slots)+len(timeOffRecords)+len(holidayInfos))
	timeOffIntervals := make([]timeInterval, 0, len(timeOffRecords))
	for _, record := range timeOffRecords {
		startLocal := record.StartDate().In(loc)
		endLocal := record.EndDate().In(loc)
		timeOffIntervals = append(timeOffIntervals, timeInterval{start: startLocal.UTC(), end: endLocal.UTC()})

		groupID := buildGroupID(startLocal, nil, AgendaEntrySourceTimeOff)
		entries = append(entries, AgendaSlot{
			PhotographerID: input.PhotographerID,
			Start:          startLocal,
			End:            endLocal,
			Status:         photosessionmodel.SlotStatusBlocked,
			GroupID:        groupID,
			Source:         AgendaEntrySourceTimeOff,
			IsTimeOff:      true,
			Timezone:       resultTimezone,
		})
	}

	for _, slot := range slots {
		period := slot.Period()
		entryPeriod := period
		startLocal := slot.SlotStart().In(loc)
		endLocal := slot.SlotEnd().In(loc)
		groupID := buildGroupID(startLocal, &entryPeriod, AgendaEntrySourceSlot)
		agendaSlot := AgendaSlot{
			SlotID:         slot.ID(),
			PhotographerID: slot.PhotographerUserID(),
			Start:          startLocal,
			End:            endLocal,
			Period:         &entryPeriod,
			Status:         slot.Status(),
			GroupID:        groupID,
			Source:         AgendaEntrySourceSlot,
			Timezone:       resultTimezone,
		}

		dateKey := dateKeyFromTime(startLocal)
		if infos, ok := holidayByDate[dateKey]; ok {
			agendaSlot.Status = photosessionmodel.SlotStatusBlocked
			agendaSlot.IsHoliday = true
			agendaSlot.HolidayCalendarIDs = collectHolidayIDs(infos)
			agendaSlot.HolidayLabels = collectHolidayLabels(infos)
			matchedHolidayDates[dateKey] = struct{}{}
		}

		for _, interval := range timeOffIntervals {
			if intervalsOverlap(slot.SlotStart(), slot.SlotEnd(), interval.start, interval.end) {
				agendaSlot.Status = photosessionmodel.SlotStatusBlocked
				agendaSlot.IsTimeOff = true
				break
			}
		}

		entries = append(entries, agendaSlot)
	}

	for key, infos := range holidayByDate {
		if _, matched := matchedHolidayDates[key]; matched {
			continue
		}
		info := infos[0]
		day := info.date.In(loc)
		startLocal := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, loc)
		endLocal := startLocal.Add(24 * time.Hour)
		groupID := buildGroupID(day, nil, AgendaEntrySourceHoliday)
		entries = append(entries, AgendaSlot{
			PhotographerID:     input.PhotographerID,
			Start:              startLocal,
			End:                endLocal,
			Status:             photosessionmodel.SlotStatusBlocked,
			GroupID:            groupID,
			Source:             AgendaEntrySourceHoliday,
			IsHoliday:          true,
			HolidayCalendarIDs: collectHolidayIDs(infos),
			HolidayLabels:      collectHolidayLabels(infos),
			Timezone:           resultTimezone,
		})
	}

	return ListAgendaOutput{Slots: entries, Total: total, Page: page, Size: size, Timezone: resultTimezone}, nil
}

func (s *photoSessionService) ensurePhotographerAgendaWithPrepared(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput, prepared *preparedEnsureContext) error {
	if tx == nil {
		return utils.InternalError("")
	}
	if prepared == nil {
		var err error
		prepared, err = s.prepareEnsureContext(input)
		if err != nil {
			return err
		}
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err = s.ensureSlots(ctx, tx, input, prepared); err != nil {
		return err
	}

	if _, err = s.repo.DeleteSlotsOutsideRange(ctx, tx, input.PhotographerID, prepared.windowStartUTC, prepared.windowEndUTC); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.ensure_agenda.cleanup_error", "err", err, "photographer_id", input.PhotographerID)
		return utils.InternalError("")
	}

	return nil
}

func (s *photoSessionService) ensureSlots(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput, prepared *preparedEnsureContext) error {
	existingSlots, err := s.repo.ListSlotsByRange(ctx, tx, input.PhotographerID, prepared.windowStartUTC, prepared.windowEndUTC)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("photo_session.ensure_agenda.list_existing_error", "err", err, "photographer_id", input.PhotographerID)
		return utils.InternalError("")
	}

	existingByStart := make(map[time.Time]photosessionmodel.PhotographerSlotInterface, len(existingSlots))
	for _, slot := range existingSlots {
		existingByStart[slot.SlotStart().UTC()] = slot
	}

	timeOffEntries, err := s.repo.ListTimeOff(ctx, tx, input.PhotographerID, prepared.windowStartUTC, prepared.windowEndUTC)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("photo_session.ensure_agenda.list_time_off_error", "err", err, "photographer_id", input.PhotographerID)
		return utils.InternalError("")
	}

	blockedDays := prepared.blockedDays
	for _, entry := range timeOffEntries {
		markDateRange(blockedDays, entry.StartDate(), entry.EndDate(), prepared.location)
	}

	nowUTC := s.now().UTC()
	newSlots := make([]photosessionmodel.PhotographerSlotInterface, 0)

	defaults, err := s.resolveDefaultAvailability(ctx, tx, input.PhotographerID)
	if err != nil {
		return err
	}

	defaultsByWeekday := make(map[time.Weekday][]photosessionmodel.PhotographerDefaultAvailabilityInterface, len(defaults))
	for _, record := range defaults {
		defaultsByWeekday[record.Weekday()] = append(defaultsByWeekday[record.Weekday()], record)
	}

	for day := prepared.windowStartLocal; day.Before(prepared.windowEndLocal); day = day.AddDate(0, 0, 1) {
		if day.Weekday() == time.Saturday || day.Weekday() == time.Sunday {
			continue
		}

		if _, isBlocked := blockedDays[dateKeyFromTime(day)]; isBlocked {
			continue
		}

		records := defaultsByWeekday[day.Weekday()]
		if len(records) == 0 {
			continue
		}

		for _, record := range records {
			slotDuration := record.SlotDurationMinutes()
			if slotDuration <= 0 {
				slotDuration = s.cfg.SlotDurationMinutes
			}
			periodStart := record.StartHour()
			if periodStart <= 0 {
				if record.Period() == photosessionmodel.SlotPeriodAfternoon {
					periodStart = s.cfg.AfternoonStartHour
				} else {
					periodStart = s.cfg.MorningStartHour
				}
			}

			slotsPerPeriod := record.SlotsPerPeriod()
			if slotsPerPeriod <= 0 {
				slotsPerPeriod = s.cfg.SlotsPerPeriod
			}

			for index := 0; index < slotsPerPeriod; index++ {
				baseStart := time.Date(day.Year(), day.Month(), day.Day(), periodStart, 0, 0, 0, prepared.location)
				slotStartLocal := baseStart.Add(time.Duration(index*slotDuration) * time.Minute)
				slotEndLocal := slotStartLocal.Add(time.Duration(slotDuration) * time.Minute)
				slotStartUTC := slotStartLocal.UTC()
				slotEndUTC := slotEndLocal.UTC()

				if !slotEndUTC.After(nowUTC) {
					continue
				}

				if _, exists := existingByStart[slotStartUTC]; exists {
					continue
				}

				slot := photosessionmodel.NewPhotographerSlot()
				slot.SetPhotographerUserID(input.PhotographerID)
				slot.SetSlotDate(slotStartUTC.Truncate(24 * time.Hour))
				slot.SetSlotStart(slotStartUTC)
				slot.SetSlotEnd(slotEndUTC)
				slot.SetStatus(photosessionmodel.SlotStatusAvailable)
				slot.SetPeriod(record.Period())

				newSlots = append(newSlots, slot)
			}
		}
	}

	if len(newSlots) == 0 {
		return nil
	}

	if err = s.repo.BulkUpsertSlots(ctx, tx, newSlots); err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("photo_session.ensure_agenda.bulk_upsert_error", "err", err)
		return utils.InternalError("")
	}

	return nil
}

func (s *photoSessionService) resolveDefaultAvailability(ctx context.Context, tx *sql.Tx, photographerID uint64) ([]photosessionmodel.PhotographerDefaultAvailabilityInterface, error) {
	records, err := s.repo.ListDefaultAvailability(ctx, tx, photographerID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("photo_session.ensure_agenda.list_default_availability_error", "err", err, "photographer_id", photographerID)
		return nil, utils.InternalError("")
	}
	if len(records) > 0 {
		return records, nil
	}
	fallback := s.buildFallbackAvailability(photographerID)
	if err = s.repo.ReplaceDefaultAvailability(ctx, tx, photographerID, fallback); err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("photo_session.ensure_agenda.persist_default_availability_error", "err", err, "photographer_id", photographerID)
		return nil, utils.InternalError("")
	}
	return fallback, nil
}

func (s *photoSessionService) buildFallbackAvailability(photographerID uint64) []photosessionmodel.PhotographerDefaultAvailabilityInterface {
	result := make([]photosessionmodel.PhotographerDefaultAvailabilityInterface, 0, 10)
	periods := []struct {
		period photosessionmodel.SlotPeriod
		start  int
	}{
		{photosessionmodel.SlotPeriodMorning, s.cfg.MorningStartHour},
		{photosessionmodel.SlotPeriodAfternoon, s.cfg.AfternoonStartHour},
	}

	for weekday := time.Monday; weekday <= time.Friday; weekday++ {
		for _, p := range periods {
			record := photosessionmodel.NewPhotographerDefaultAvailability()
			record.SetPhotographerUserID(photographerID)
			record.SetWeekday(weekday)
			record.SetPeriod(p.period)
			record.SetStartHour(p.start)
			record.SetSlotsPerPeriod(s.cfg.SlotsPerPeriod)
			record.SetSlotDurationMinutes(s.cfg.SlotDurationMinutes)
			result = append(result, record)
		}
	}

	return result
}

type holidayDateInfo struct {
	calendarID uint64
	date       time.Time
	label      string
}

type timeInterval struct {
	start time.Time
	end   time.Time
}

func (s *photoSessionService) fetchHolidayDates(ctx context.Context, calendarIDs []uint64, from, to time.Time) ([]holidayDateInfo, error) {
	if len(calendarIDs) == 0 {
		return nil, nil
	}

	fromPtr := from
	toPtr := to
	if toPtr.Before(fromPtr) {
		toPtr = fromPtr
	}

	entries := make([]holidayDateInfo, 0)

	for _, calendarID := range calendarIDs {
		if calendarID == 0 {
			continue
		}

		filter := holidaymodel.CalendarDatesFilter{
			CalendarID: calendarID,
			From:       &fromPtr,
			To:         &toPtr,
			Limit:      200,
			Page:       1,
		}

		for {
			result, err := s.holidayService.ListCalendarDates(ctx, filter)
			if err != nil {
				utils.SetSpanError(ctx, err)
				utils.LoggerFromContext(ctx).Error("service.list_agenda.holiday_error", "err", err, "calendar_id", calendarID)
				return nil, derrors.Wrap(err, derrors.KindInfra, "failed to list holiday dates")
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

	return entries, nil
}

func collectHolidayIDs(infos []holidayDateInfo) []uint64 {
	if len(infos) == 0 {
		return nil
	}
	seen := make(map[uint64]struct{}, len(infos))
	for _, info := range infos {
		seen[info.calendarID] = struct{}{}
	}
	ids := make([]uint64, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

func collectHolidayLabels(infos []holidayDateInfo) []string {
	if len(infos) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(infos))
	labels := make([]string, 0, len(infos))
	for _, info := range infos {
		label := strings.TrimSpace(info.label)
		if label == "" {
			label = "Holiday"
		}
		if _, ok := seen[label]; ok {
			continue
		}
		seen[label] = struct{}{}
		labels = append(labels, label)
	}
	sort.Strings(labels)
	return labels
}

func intervalsOverlap(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}

func buildGroupID(day time.Time, period *photosessionmodel.SlotPeriod, source AgendaEntrySource) string {
	key := day.Format("2006-01-02")
	if period != nil {
		key = fmt.Sprintf("%s-%s", key, string(*period))
	}
	return fmt.Sprintf("%s-%s", strings.ToLower(string(source)), key)
}

// func (s *photoSessionService) prepareEnsureContext(ctx context.Context, input EnsureAgendaInput) (*preparedEnsureContext, error) {
func (s *photoSessionService) prepareEnsureContext(input EnsureAgendaInput) (*preparedEnsureContext, error) {
	if input.PhotographerID == 0 {
		return nil, utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}

	loc, err := resolveLocation(input.Timezone)
	if err != nil {
		return nil, err
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
		return nil, utils.ValidationError("workdayEndHour", "workdayEndHour must be greater than workdayStartHour")
	}

	horizonMonths := input.HorizonMonths
	if horizonMonths <= 0 {
		horizonMonths = defaultHorizonMonths
	}

	now := s.now().In(loc)
	windowStartLocal := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	windowEndLocal := windowStartLocal.AddDate(0, horizonMonths, 0)

	return &preparedEnsureContext{
		location:         loc,
		windowStartLocal: windowStartLocal,
		windowEndLocal:   windowEndLocal,
		windowStartUTC:   windowStartLocal.UTC(),
		windowEndUTC:     windowEndLocal.UTC(),
		workdayStartHour: workdayStart,
		workdayEndHour:   workdayEnd,
		blockedDays:      make(map[string]struct{}),
	}, nil
}

// func (s *photoSessionService) loadHolidayDays(ctx context.Context, input EnsureAgendaInput, from, to time.Time) (map[string]struct{}, error) {
// 	blocked := make(map[string]struct{})
// 	if input.HolidayCalendarID == nil || *input.HolidayCalendarID == 0 {
// 		return blocked, nil
// 	}

// 	filter := holidaymodel.CalendarDatesFilter{
// 		CalendarID: *input.HolidayCalendarID,
// 		From:       &from,
// 		To:         &to,
// 		Limit:      100,
// 		Page:       1,
// 	}

// 	result, err := s.holidayService.ListCalendarDates(ctx, filter)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, date := range result.Dates {
// 		blocked[dateKeyFromTime(date.HolidayDate().In(from.Location()))] = struct{}{}
// 	}

// 	return blocked, nil
// }

func markDateRange(target map[string]struct{}, start, end time.Time, loc *time.Location) {
	start = start.In(loc)
	end = end.In(loc)

	for d := start; d.Before(end) || d.Equal(end); d = d.AddDate(0, 0, 1) {
		target[dateKeyFromTime(d)] = struct{}{}
	}
}

func dateKeyFromTime(t time.Time) string {
	return t.Format("2006-01-02")
}

func resolveLocation(timezone string) (*time.Location, error) {
	if timezone == "" {
		timezone = defaultTimezone
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, utils.ValidationError("timezone", "Invalid timezone")
	}

	return loc, nil
}

// UpdateSessionStatusInput contains the required data to update a photo session status.
type UpdateSessionStatusInput struct {
	SessionID      uint64
	PhotographerID uint64
	Status         string
}

// EnsureAgendaInput controls agenda generation parameters.
type EnsureAgendaInput struct {
	PhotographerID    uint64
	Timezone          string
	HolidayCalendarID *uint64
	HorizonMonths     int
	WorkdayStartHour  int
	WorkdayEndHour    int
}

// TimeOffInput represents the payload to block a photographer agenda.
type TimeOffInput struct {
	PhotographerID    uint64
	StartDate         time.Time
	EndDate           time.Time
	Reason            *string
	Timezone          string
	HolidayCalendarID *uint64
	HorizonMonths     int
	WorkdayStartHour  int
	WorkdayEndHour    int
}

func (input TimeOffInput) toModel() photosessionmodel.PhotographerTimeOffInterface {
	model := photosessionmodel.NewPhotographerTimeOff()
	model.SetPhotographerUserID(input.PhotographerID)
	model.SetStartDate(input.StartDate.UTC())
	model.SetEndDate(input.EndDate.UTC())
	if input.Reason != nil {
		reason := strings.TrimSpace(*input.Reason)
		if reason != "" {
			reasonCopy := reason
			model.SetReason(&reasonCopy)
		}
	}
	return model
}

// DeleteTimeOffInput represents the payload to unblock a photographer agenda.
type DeleteTimeOffInput struct {
	TimeOffID         uint64
	PhotographerID    uint64
	Timezone          string
	HolidayCalendarID *uint64
	HorizonMonths     int
	WorkdayStartHour  int
	WorkdayEndHour    int
}

func validateTimeOffInput(input TimeOffInput) error {
	if input.PhotographerID == 0 {
		return utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}
	if input.StartDate.IsZero() {
		return utils.ValidationError("startDate", "startDate is required")
	}
	if input.EndDate.IsZero() {
		return utils.ValidationError("endDate", "endDate is required")
	}
	if input.EndDate.Before(input.StartDate) {
		return utils.ValidationError("endDate", "endDate must be greater than or equal to startDate")
	}
	if _, err := resolveLocation(input.Timezone); err != nil {
		return err
	}
	if input.Reason != nil {
		reason := strings.TrimSpace(*input.Reason)
		if len(reason) > maxTimeOffReasonLength {
			return utils.ValidationError("reason", fmt.Sprintf("reason must be at most %d characters", maxTimeOffReasonLength))
		}
	}
	return nil
}

func validateListAgendaInput(input ListAgendaInput) error {
	if input.PhotographerID == 0 {
		return derrors.Validation("photographerId must be greater than zero", nil)
	}
	if input.StartDate.IsZero() {
		return derrors.Validation("startDate is required", nil)
	}
	if input.EndDate.IsZero() {
		return derrors.Validation("endDate is required", nil)
	}
	if input.EndDate.Before(input.StartDate) {
		return derrors.Validation("endDate must be after or equal to startDate", nil)
	}
	if input.Size < 0 {
		return derrors.Validation("size must be zero or greater", nil)
	}
	if input.Page < 0 {
		return derrors.Validation("page must be zero or greater", nil)
	}
	if input.Timezone != "" {
		if _, err := resolveLocation(input.Timezone); err != nil {
			return err
		}
	}
	return nil
}

type preparedEnsureContext struct {
	location         *time.Location
	windowStartLocal time.Time
	windowEndLocal   time.Time
	windowStartUTC   time.Time
	windowEndUTC     time.Time
	workdayStartHour int
	workdayEndHour   int
	blockedDays      map[string]struct{}
}

// ListAgendaInput defines the input for listing a photographer's agenda.
type ListAgendaInput struct {
	PhotographerID     uint64
	StartDate          time.Time
	EndDate            time.Time
	Page               int
	Size               int
	Timezone           string
	HolidayCalendarIDs []uint64
}

// ListAgendaOutput defines the output for a photographer's agenda.
type ListAgendaOutput struct {
	Slots    []AgendaSlot `json:"slots"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	Size     int          `json:"size"`
	Timezone string       `json:"timezone"`
}

// AgendaSlot is a DTO representing a consolidated agenda entry.
type AgendaSlot struct {
	SlotID             uint64                        `json:"slotId,omitempty"`
	PhotographerID     uint64                        `json:"photographerId,omitempty"`
	Start              time.Time                     `json:"start"`
	End                time.Time                     `json:"end"`
	Period             *photosessionmodel.SlotPeriod `json:"period,omitempty"`
	Status             photosessionmodel.SlotStatus  `json:"status"`
	GroupID            string                        `json:"groupId"`
	Source             AgendaEntrySource             `json:"source"`
	IsHoliday          bool                          `json:"isHoliday"`
	IsTimeOff          bool                          `json:"isTimeOff"`
	HolidayLabels      []string                      `json:"holidayLabels,omitempty"`
	HolidayCalendarIDs []uint64                      `json:"holidayCalendarIds,omitempty"`
	Timezone           string                        `json:"timezone"`
}

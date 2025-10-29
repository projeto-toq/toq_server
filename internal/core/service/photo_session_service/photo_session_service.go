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
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
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

// PhotoSessionServiceInterface exposes orchestration helpers around photographer agenda entries.
type PhotoSessionServiceInterface interface {
	EnsurePhotographerAgenda(ctx context.Context, input EnsureAgendaInput) error
	EnsurePhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error
	RefreshPhotographerAgenda(ctx context.Context, input EnsureAgendaInput) error
	RefreshPhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error
	CreateTimeOff(ctx context.Context, input TimeOffInput) (uint64, error)
	CreateTimeOffWithTx(ctx context.Context, tx *sql.Tx, input TimeOffInput) (uint64, error)
	DeleteTimeOff(ctx context.Context, input DeleteTimeOffInput) error
	DeleteTimeOffWithTx(ctx context.Context, tx *sql.Tx, input DeleteTimeOffInput) error
	ListTimeOff(ctx context.Context, input ListTimeOffInput) (ListTimeOffOutput, error)
	GetTimeOffDetail(ctx context.Context, input TimeOffDetailInput) (TimeOffDetailResult, error)
	UpdateTimeOff(ctx context.Context, input UpdateTimeOffInput) (TimeOffDetailResult, error)
	UpdateSessionStatus(ctx context.Context, input UpdateSessionStatusInput) error
	ListAgenda(ctx context.Context, input ListAgendaInput) (ListAgendaOutput, error)
	ListAvailability(ctx context.Context, input ListAvailabilityInput) (ListAvailabilityOutput, error)
	ReservePhotoSession(ctx context.Context, input ReserveSessionInput) (ReserveSessionOutput, error)
	ConfirmPhotoSession(ctx context.Context, input ConfirmSessionInput) (ConfirmSessionOutput, error)
	CancelPhotoSession(ctx context.Context, input CancelSessionInput) (CancelSessionOutput, error)
}

type photoSessionService struct {
	repo           photosessionrepository.PhotoSessionRepositoryInterface
	holidayRepo    photosessionrepository.PhotographerHolidayCalendarRepository
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
	var holidayRepo photosessionrepository.PhotographerHolidayCalendarRepository
	if value, ok := repo.(photosessionrepository.PhotographerHolidayCalendarRepository); ok {
		holidayRepo = value
	}

	return &photoSessionService{
		repo:           repo,
		holidayRepo:    holidayRepo,
		listingRepo:    listingRepo,
		holidayService: holidayService,
		globalService:  globalService,
		cfg:            normalizeConfig(cfg),
		now:            time.Now,
	}
}

// EnsurePhotographerAgenda provisions bootstrap agenda entries for a photographer.
func (s *photoSessionService) EnsurePhotographerAgenda(ctx context.Context, input EnsureAgendaInput) error {
	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = s.globalService.RollbackTransaction(ctx, tx)
		}
	}()

	if err := s.ensurePhotographerAgendaInternal(ctx, tx, input); err != nil {
		return err
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		return err
	}
	committed = true
	return nil
}

// EnsurePhotographerAgendaWithTx provisions bootstrap agenda entries using an existing transaction.
func (s *photoSessionService) EnsurePhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error {
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
	if tx == nil {
		return utils.InternalError("")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
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
			logger.Error("photo_session.bootstrap.replace_holiday_assoc_error", "photographer_id", input.PhotographerID, "err", err)
			return utils.InternalError("")
		}
	}

	if _, err := s.repo.DeleteEntriesBySource(ctx, tx, input.PhotographerID, photosessionmodel.AgendaEntryTypeBlock, photosessionmodel.AgendaEntrySourceOnboarding, nil); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.bootstrap.delete_block_entries_error", "photographer_id", input.PhotographerID, "err", err)
		return utils.InternalError("")
	}

	if _, err := s.repo.DeleteEntriesBySource(ctx, tx, input.PhotographerID, photosessionmodel.AgendaEntryTypeHoliday, photosessionmodel.AgendaEntrySourceHoliday, nil); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.bootstrap.delete_holiday_entries_error", "photographer_id", input.PhotographerID, "err", err)
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
		logger.Error("photo_session.bootstrap.create_entries_error", "photographer_id", input.PhotographerID, "err", err)
		return utils.InternalError("")
	}

	return nil
}

// CreateTimeOff registers a new time-off entry.
func (s *photoSessionService) CreateTimeOff(ctx context.Context, input TimeOffInput) (uint64, error) {
	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		return 0, err
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

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		return 0, err
	}
	committed = true
	return id, nil
}

// CreateTimeOffWithTx registers a new time-off entry using an existing transaction.
func (s *photoSessionService) CreateTimeOffWithTx(ctx context.Context, tx *sql.Tx, input TimeOffInput) (uint64, error) {
	if tx == nil {
		return 0, utils.InternalError("")
	}
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

	loc, tzErr := resolveLocation(input.Timezone)
	if tzErr != nil {
		return 0, tzErr
	}

	entry := photosessionmodel.NewAgendaEntry()
	entry.SetPhotographerUserID(input.PhotographerID)
	entry.SetEntryType(photosessionmodel.AgendaEntryTypeTimeOff)
	entry.SetSource(photosessionmodel.AgendaEntrySourceManual)
	entry.SetStartsAt(input.StartDate.UTC())
	entry.SetEndsAt(input.EndDate.UTC())
	entry.SetBlocking(true)
	entry.SetTimezone(loc.String())
	if input.Reason != nil {
		reason := strings.TrimSpace(*input.Reason)
		if reason != "" {
			entry.SetReason(reason)
		}
	}

	ids, err := s.repo.CreateEntries(ctx, tx, []photosessionmodel.AgendaEntryInterface{entry})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.time_off.create_error", "photographer_id", input.PhotographerID, "err", err)
		return 0, utils.InternalError("")
	}
	if len(ids) == 0 {
		return 0, utils.InternalError("")
	}

	return ids[0], nil
}

// DeleteTimeOff removes an existing time-off entry.
func (s *photoSessionService) DeleteTimeOff(ctx context.Context, input DeleteTimeOffInput) error {
	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		return err
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

// DeleteTimeOffWithTx removes an existing time-off entry inside a transaction.
func (s *photoSessionService) DeleteTimeOffWithTx(ctx context.Context, tx *sql.Tx, input DeleteTimeOffInput) error {
	if tx == nil {
		return utils.InternalError("")
	}
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

	entry, err := s.repo.GetEntryByIDForUpdate(ctx, tx, input.TimeOffID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Time off")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.time_off.get_error", "time_off_id", input.TimeOffID, "err", err)
		return utils.InternalError("")
	}

	if entry.EntryType() != photosessionmodel.AgendaEntryTypeTimeOff || entry.PhotographerUserID() != input.PhotographerID {
		return utils.NotFoundError("Time off")
	}

	if err := s.repo.DeleteEntryByID(ctx, tx, input.TimeOffID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.time_off.delete_error", "time_off_id", input.TimeOffID, "err", err)
		return utils.InternalError("")
	}

	return nil
}

// ListTimeOff returns paginated agenda entries of type TIME_OFF.
func (s *photoSessionService) ListTimeOff(ctx context.Context, input ListTimeOffInput) (ListTimeOffOutput, error) {
	if input.PhotographerID == 0 {
		return ListTimeOffOutput{}, utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}
	if input.RangeFrom.IsZero() {
		return ListTimeOffOutput{}, utils.ValidationError("rangeFrom", "rangeFrom is required")
	}
	if input.RangeTo.IsZero() {
		return ListTimeOffOutput{}, utils.ValidationError("rangeTo", "rangeTo is required")
	}
	if input.RangeTo.Before(input.RangeFrom) {
		return ListTimeOffOutput{}, utils.ValidationError("rangeTo", "rangeTo must be greater than or equal to rangeFrom")
	}

	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "service.ListTimeOff")
	if err != nil {
		return ListTimeOffOutput{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	loc, tzErr := resolveLocation(input.Timezone)
	if tzErr != nil {
		return ListTimeOffOutput{}, tzErr
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

	tx, err := s.globalService.StartReadOnlyTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.list_time_off.tx_start_error", "err", err)
		return ListTimeOffOutput{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("photo_session.list_time_off.tx_rollback_error", "err", rbErr)
		}
	}()

	entries, err := s.repo.ListEntriesByRange(ctx, tx, input.PhotographerID, input.RangeFrom.UTC(), input.RangeTo.UTC())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.list_time_off.repo_error", "photographer_id", input.PhotographerID, "err", err)
		return ListTimeOffOutput{}, utils.InternalError("")
	}

	timeOffEntries := make([]photosessionmodel.AgendaEntryInterface, 0)
	for _, entry := range entries {
		if entry.EntryType() != photosessionmodel.AgendaEntryTypeTimeOff {
			continue
		}
		entry.SetStartsAt(utils.ConvertToLocation(entry.StartsAt(), loc))
		entry.SetEndsAt(utils.ConvertToLocation(entry.EndsAt(), loc))
		timeOffEntries = append(timeOffEntries, entry)
	}

	sort.Slice(timeOffEntries, func(i, j int) bool {
		if timeOffEntries[i].StartsAt().Equal(timeOffEntries[j].StartsAt()) {
			return timeOffEntries[i].ID() < timeOffEntries[j].ID()
		}
		return timeOffEntries[i].StartsAt().Before(timeOffEntries[j].StartsAt())
	})

	total := len(timeOffEntries)
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}

	return ListTimeOffOutput{
		TimeOffs: timeOffEntries[start:end],
		Total:    int64(total),
		Page:     page,
		Size:     size,
		Timezone: loc.String(),
	}, nil
}

// GetTimeOffDetail fetches a specific time-off entry ensuring ownership.
func (s *photoSessionService) GetTimeOffDetail(ctx context.Context, input TimeOffDetailInput) (TimeOffDetailResult, error) {
	if input.TimeOffID == 0 {
		return TimeOffDetailResult{}, utils.ValidationError("timeOffId", "timeOffId must be greater than zero")
	}
	if input.PhotographerID == 0 {
		return TimeOffDetailResult{}, utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "service.GetTimeOffDetail")
	if err != nil {
		return TimeOffDetailResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	loc, tzErr := resolveLocation(input.Timezone)
	if tzErr != nil {
		return TimeOffDetailResult{}, tzErr
	}

	tx, err := s.globalService.StartReadOnlyTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.get_time_off.tx_start_error", "err", err)
		return TimeOffDetailResult{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("photo_session.get_time_off.tx_rollback_error", "err", rbErr)
		}
	}()

	entry, err := s.repo.GetEntryByID(ctx, tx, input.TimeOffID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TimeOffDetailResult{}, utils.NotFoundError("Time off")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.get_time_off.repo_error", "time_off_id", input.TimeOffID, "err", err)
		return TimeOffDetailResult{}, utils.InternalError("")
	}

	if entry.EntryType() != photosessionmodel.AgendaEntryTypeTimeOff || entry.PhotographerUserID() != input.PhotographerID {
		return TimeOffDetailResult{}, utils.NotFoundError("Time off")
	}

	entry.SetStartsAt(utils.ConvertToLocation(entry.StartsAt(), loc))
	entry.SetEndsAt(utils.ConvertToLocation(entry.EndsAt(), loc))

	return TimeOffDetailResult{TimeOff: entry, Timezone: loc.String()}, nil
}

// UpdateTimeOff updates a time-off entry and returns the mutated record.
func (s *photoSessionService) UpdateTimeOff(ctx context.Context, input UpdateTimeOffInput) (TimeOffDetailResult, error) {
	if input.TimeOffID == 0 {
		return TimeOffDetailResult{}, utils.ValidationError("timeOffId", "timeOffId must be greater than zero")
	}
	if err := validateTimeOffInput(TimeOffInput{
		PhotographerID:    input.PhotographerID,
		StartDate:         input.StartDate,
		EndDate:           input.EndDate,
		Reason:            input.Reason,
		Timezone:          input.Timezone,
		HolidayCalendarID: input.HolidayCalendarID,
		HorizonMonths:     input.HorizonMonths,
		WorkdayStartHour:  input.WorkdayStartHour,
		WorkdayEndHour:    input.WorkdayEndHour,
	}); err != nil {
		return TimeOffDetailResult{}, err
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return TimeOffDetailResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	loc, tzErr := resolveLocation(input.Timezone)
	if tzErr != nil {
		return TimeOffDetailResult{}, tzErr
	}

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_time_off.tx_start_error", "err", err)
		return TimeOffDetailResult{}, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("photo_session.update_time_off.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	entry, err := s.repo.GetEntryByIDForUpdate(ctx, tx, input.TimeOffID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TimeOffDetailResult{}, utils.NotFoundError("Time off")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_time_off.get_error", "time_off_id", input.TimeOffID, "err", err)
		return TimeOffDetailResult{}, utils.InternalError("")
	}

	if entry.EntryType() != photosessionmodel.AgendaEntryTypeTimeOff || entry.PhotographerUserID() != input.PhotographerID {
		return TimeOffDetailResult{}, utils.NotFoundError("Time off")
	}

	entry.SetStartsAt(input.StartDate.UTC())
	entry.SetEndsAt(input.EndDate.UTC())
	entry.SetTimezone(loc.String())
	entry.SetSource(photosessionmodel.AgendaEntrySourceManual)
	entry.SetBlocking(true)
	entry.ClearReason()
	if input.Reason != nil {
		reason := strings.TrimSpace(*input.Reason)
		if reason != "" {
			entry.SetReason(reason)
		}
	}

	if err := s.repo.UpdateEntry(ctx, tx, entry); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_time_off.repo_error", "time_off_id", input.TimeOffID, "err", err)
		return TimeOffDetailResult{}, utils.InternalError("")
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_time_off.tx_commit_error", "err", err)
		return TimeOffDetailResult{}, utils.InternalError("")
	}
	committed = true

	entry.SetStartsAt(utils.ConvertToLocation(entry.StartsAt(), loc))
	entry.SetEndsAt(utils.ConvertToLocation(entry.EndsAt(), loc))

	return TimeOffDetailResult{TimeOff: entry, Timezone: loc.String()}, nil
}

// UpdateSessionStatus updates the status of a photo session booking.
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

	booking, err := s.repo.GetBookingByIDForUpdate(ctx, tx, input.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("session not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.get_booking_error", "session_id", input.SessionID, "err", err)
		return derrors.Infra("failed to load session booking", err)
	}

	if booking.PhotographerUserID() != input.PhotographerID {
		return derrors.Forbidden("session does not belong to photographer")
	}
	if booking.Status() != photosessionmodel.BookingStatusPendingApproval {
		return derrors.Conflict("session is not pending approval")
	}

	_, err = s.repo.GetEntryByIDForUpdate(ctx, tx, booking.AgendaEntryID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.Infra("agenda entry missing for booking", err)
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.get_entry_error", "agenda_entry_id", booking.AgendaEntryID(), "err", err)
		return derrors.Infra("failed to load agenda entry", err)
	}

	if err := s.repo.UpdateBookingStatus(ctx, tx, booking.ID(), status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("session not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.update_error", "session_id", booking.ID(), "err", err)
		return derrors.Infra("failed to update session status", err)
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_status.tx_commit_error", "session_id", booking.ID(), "err", err)
		return derrors.Infra("failed to commit transaction", err)
	}
	committed = true

	logger.Info("photo_session.status.updated", "session_id", booking.ID(), "photographer_id", input.PhotographerID, "status", statusStr)
	return nil
}

// ListAgenda returns agenda entries within the interval.
func (s *photoSessionService) ListAgenda(ctx context.Context, input ListAgendaInput) (ListAgendaOutput, error) {
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "service.ListAgenda")
	if err != nil {
		return ListAgendaOutput{}, derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	if err := validateListAgendaInput(input); err != nil {
		return ListAgendaOutput{}, err
	}

	loc, tzErr := resolveLocation(input.Timezone)
	if tzErr != nil {
		return ListAgendaOutput{}, tzErr
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

	tx, err := s.globalService.StartReadOnlyTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("service.list_agenda.tx_start_error", "err", err)
		return ListAgendaOutput{}, derrors.Wrap(err, derrors.KindInfra, "failed to start transaction")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			utils.LoggerFromContext(ctx).Error("service.list_agenda.tx_rollback_error", "err", rbErr)
		}
	}()

	entries, err := s.repo.ListEntriesByRange(ctx, tx, input.PhotographerID, input.StartDate.UTC(), input.EndDate.UTC())
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("service.list_agenda.repo_error", "photographer_id", input.PhotographerID, "err", err)
		return ListAgendaOutput{}, derrors.Wrap(err, derrors.KindInfra, "failed to list agenda entries")
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].StartsAt().Equal(entries[j].StartsAt()) {
			return entries[i].ID() < entries[j].ID()
		}
		return entries[i].StartsAt().Before(entries[j].StartsAt())
	})

	slots := make([]AgendaSlot, 0, len(entries))
	for _, entry := range entries {
		slots = append(slots, s.buildAgendaSlot(entry, loc))
	}

	total := len(slots)
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}

	return ListAgendaOutput{
		Slots:    slots[start:end],
		Total:    int64(total),
		Page:     page,
		Size:     size,
		Timezone: loc.String(),
	}, nil
}

// ListAvailability computes booking availability windows for photographers.
func (s *photoSessionService) ListAvailability(ctx context.Context, input ListAvailabilityInput) (ListAvailabilityOutput, error) {
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "service.ListAvailability")
	if err != nil {
		return ListAvailabilityOutput{}, derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	loc, tzErr := resolveLocation(input.Timezone)
	if tzErr != nil {
		return ListAvailabilityOutput{}, tzErr
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

	workdayStart := input.WorkdayStartHour
	if workdayStart <= 0 {
		workdayStart = defaultWorkdayStartHour
	}
	workdayEnd := input.WorkdayEndHour
	if workdayEnd <= 0 {
		workdayEnd = defaultWorkdayEndHour
	}
	if workdayEnd <= workdayStart {
		return ListAvailabilityOutput{}, derrors.Validation("workdayEndHour must be greater than workdayStartHour", nil)
	}

	now := s.now().In(loc)
	var rangeStart time.Time
	if input.From != nil {
		rangeStart = input.From.In(loc)
	} else {
		rangeStart = now
	}
	var rangeEnd time.Time
	if input.To != nil {
		rangeEnd = input.To.In(loc)
	} else {
		rangeEnd = rangeStart.AddDate(0, defaultHorizonMonths, 0)
	}
	if rangeEnd.Before(rangeStart) {
		return ListAvailabilityOutput{}, derrors.Validation("to must be after from", nil)
	}

	slotDuration := time.Duration(s.cfg.SlotDurationMinutes) * time.Minute
	if slotDuration <= 0 {
		slotDuration = 4 * time.Hour
	}

	filterPeriod := input.Period

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("photo_session.list_availability.tx_start_error", "err", txErr)
		return ListAvailabilityOutput{}, derrors.Infra("failed to start transaction", txErr)
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("photo_session.list_availability.tx_rollback_error", "err", rbErr)
		}
	}()

	photographerIDs, repoErr := s.repo.ListPhotographerIDs(ctx, tx)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("photo_session.list_availability.list_photographers_error", "err", repoErr)
		return ListAvailabilityOutput{}, derrors.Infra("failed to list photographers", repoErr)
	}

	availability := make([]AvailabilitySlot, 0)
	for _, photographerID := range photographerIDs {
		entries, err := s.repo.ListEntriesByRange(ctx, tx, photographerID, rangeStart.UTC(), rangeEnd.UTC())
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("photo_session.list_availability.list_entries_error", "photographer_id", photographerID, "err", err)
			return ListAvailabilityOutput{}, derrors.Infra("failed to load agenda entries", err)
		}

		workingRanges := buildWorkingRanges(rangeStart, rangeEnd, loc, workdayStart, workdayEnd)
		freeRanges := applyBlockingEntries(workingRanges, entries, loc)
		freeRanges = prunePastRanges(freeRanges, now)

		for _, free := range freeRanges {
			slots := splitIntoSlots(free, slotDuration)
			for _, slot := range slots {
				period := determineSlotPeriod(slot.start)
				if filterPeriod != nil && period != *filterPeriod {
					continue
				}
				id := encodeSlotID(photographerID, slot.start.UTC())
				availability = append(availability, AvailabilitySlot{
					SlotID:         id,
					PhotographerID: photographerID,
					Start:          slot.start,
					End:            slot.end,
					Period:         period,
					SourceTimezone: loc.String(),
				})
			}
		}
	}

	sortAvailabilitySlots(availability, input.Sort)

	total := len(availability)
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}

	return ListAvailabilityOutput{
		Slots:    availability[start:end],
		Total:    int64(total),
		Page:     page,
		Size:     size,
		Timezone: loc.String(),
	}, nil
}

// ReservePhotoSession blocks a slot window for a listing owner.
func (s *photoSessionService) ReservePhotoSession(ctx context.Context, input ReserveSessionInput) (ReserveSessionOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ReserveSessionOutput{}, derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.UserID <= 0 {
		return ReserveSessionOutput{}, derrors.Auth("unauthorized")
	}
	if input.ListingID <= 0 {
		return ReserveSessionOutput{}, derrors.Validation("listingId must be greater than zero", map[string]any{"listingId": "greater_than_zero"})
	}
	if input.SlotID == 0 {
		return ReserveSessionOutput{}, derrors.Validation("slotId must be greater than zero", map[string]any{"slotId": "greater_than_zero"})
	}

	photographerID, slotStartUTC := decodeSlotID(input.SlotID)
	if photographerID == 0 {
		return ReserveSessionOutput{}, derrors.Validation("slotId is invalid", map[string]any{"slotId": "invalid"})
	}

	loc, tzErr := resolveLocation("")
	if tzErr != nil {
		return ReserveSessionOutput{}, tzErr
	}

	slotDuration := time.Duration(s.cfg.SlotDurationMinutes) * time.Minute
	if slotDuration <= 0 {
		slotDuration = 4 * time.Hour
	}

	slotStart := slotStartUTC.In(loc)
	slotEnd := slotStart.Add(slotDuration)
	if !slotEnd.After(slotStart) {
		return ReserveSessionOutput{}, derrors.Validation("slot duration must be positive", map[string]any{"slot": "invalid_duration"})
	}
	if slotEnd.Before(s.now().In(loc)) {
		return ReserveSessionOutput{}, derrors.ErrSlotUnavailable
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("photo_session.reserve.tx_start_error", "err", txErr)
		return ReserveSessionOutput{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("photo_session.reserve.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	listing, err := s.listingRepo.GetListingByID(ctx, tx, input.ListingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ReserveSessionOutput{}, utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.reserve.get_listing_error", "listing_id", input.ListingID, "err", err)
		return ReserveSessionOutput{}, derrors.Infra("failed to load listing", err)
	}

	if listing.Deleted() {
		return ReserveSessionOutput{}, utils.BadRequest("listing is not available")
	}

	if listing.UserID() != input.UserID {
		return ReserveSessionOutput{}, derrors.Auth("listing does not belong to user")
	}

	if !listingAllowsPhotoSession(listing.Status()) {
		return ReserveSessionOutput{}, derrors.ErrListingNotEligible
	}

	conflicts, err := s.repo.FindBlockingEntries(ctx, tx, photographerID, slotStart.UTC(), slotEnd.UTC())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.reserve.find_blocking_error", "photographer_id", photographerID, "err", err)
		return ReserveSessionOutput{}, derrors.Infra("failed to verify photographer agenda", err)
	}
	if len(conflicts) > 0 {
		return ReserveSessionOutput{}, derrors.ErrSlotUnavailable
	}

	agendaEntry := photosessionmodel.NewAgendaEntry()
	agendaEntry.SetPhotographerUserID(photographerID)
	agendaEntry.SetEntryType(photosessionmodel.AgendaEntryTypePhotoSession)
	agendaEntry.SetSource(photosessionmodel.AgendaEntrySourceBooking)
	agendaEntry.SetStartsAt(slotStart.UTC())
	agendaEntry.SetEndsAt(slotEnd.UTC())
	agendaEntry.SetBlocking(true)
	agendaEntry.SetTimezone(loc.String())
	if input.ListingID > 0 {
		agendaEntry.SetSourceID(uint64(input.ListingID))
	}

	entryIDs, err := s.repo.CreateEntries(ctx, tx, []photosessionmodel.AgendaEntryInterface{agendaEntry})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.reserve.create_entry_error", "photographer_id", photographerID, "err", err)
		return ReserveSessionOutput{}, derrors.Infra("failed to create agenda entry", err)
	}
	if len(entryIDs) == 0 {
		return ReserveSessionOutput{}, derrors.Infra("failed to create agenda entry", fmt.Errorf("no entry id returned"))
	}
	entryID := entryIDs[0]

	booking := photosessionmodel.NewPhotoSessionBooking()
	booking.SetAgendaEntryID(entryID)
	booking.SetPhotographerUserID(photographerID)
	booking.SetListingID(input.ListingID)
	booking.SetStartsAt(slotStart.UTC())
	booking.SetEndsAt(slotEnd.UTC())
	booking.SetStatus(photosessionmodel.BookingStatusPendingApproval)

	bookingID, err := s.repo.CreateBooking(ctx, tx, booking)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.reserve.create_booking_error", "agenda_entry_id", entryID, "err", err)
		return ReserveSessionOutput{}, derrors.Infra("failed to create booking", err)
	}

	if updateErr := s.listingRepo.UpdateListingStatus(ctx, tx, listing.ID(), listingmodel.StatusPendingAvailabilityConfirm, listing.Status()); updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return ReserveSessionOutput{}, derrors.ErrListingNotEligible
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("photo_session.reserve.update_listing_status_error", "listing_id", listing.ID(), "err", updateErr)
		return ReserveSessionOutput{}, derrors.Infra("failed to update listing status", updateErr)
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("photo_session.reserve.commit_error", "listing_id", listing.ID(), "err", commitErr)
		return ReserveSessionOutput{}, derrors.Infra("failed to commit reservation", commitErr)
	}
	committed = true

	logger.Info("photo_session.reserve.success", "listing_id", listing.ID(), "booking_id", bookingID, "photographer_id", photographerID, "slot_start", slotStart)

	return ReserveSessionOutput{
		PhotoSessionID: bookingID,
		SlotID:         input.SlotID,
		SlotStart:      slotStart,
		SlotEnd:        slotEnd,
		PhotographerID: photographerID,
		ListingID:      listing.ID(),
	}, nil
}

// ConfirmPhotoSession finalizes a pending reservation after photographer acceptance.
func (s *photoSessionService) ConfirmPhotoSession(ctx context.Context, input ConfirmSessionInput) (ConfirmSessionOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ConfirmSessionOutput{}, derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.UserID <= 0 {
		return ConfirmSessionOutput{}, derrors.Auth("unauthorized")
	}
	if input.ListingID <= 0 {
		return ConfirmSessionOutput{}, derrors.Validation("listingId must be greater than zero", map[string]any{"listingId": "greater_than_zero"})
	}
	if input.PhotoSessionID == 0 {
		return ConfirmSessionOutput{}, derrors.Validation("photoSessionId must be greater than zero", map[string]any{"photoSessionId": "greater_than_zero"})
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("photo_session.confirm.tx_start_error", "err", txErr)
		return ConfirmSessionOutput{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("photo_session.confirm.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	listing, err := s.listingRepo.GetListingByID(ctx, tx, input.ListingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ConfirmSessionOutput{}, utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.confirm.get_listing_error", "listing_id", input.ListingID, "err", err)
		return ConfirmSessionOutput{}, derrors.Infra("failed to load listing", err)
	}

	if listing.Deleted() {
		return ConfirmSessionOutput{}, utils.BadRequest("listing is not available")
	}

	if listing.UserID() != input.UserID {
		return ConfirmSessionOutput{}, derrors.Auth("listing does not belong to user")
	}

	if !listingAllowsPhotoSession(listing.Status()) {
		return ConfirmSessionOutput{}, derrors.ErrListingNotEligible
	}

	booking, err := s.repo.GetBookingByIDForUpdate(ctx, tx, input.PhotoSessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ConfirmSessionOutput{}, utils.NotFoundError("Photo session")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.confirm.get_booking_error", "photo_session_id", input.PhotoSessionID, "err", err)
		return ConfirmSessionOutput{}, derrors.Infra("failed to load booking", err)
	}

	if booking.ListingID() != input.ListingID {
		return ConfirmSessionOutput{}, derrors.Auth("photo session does not belong to listing")
	}

	switch booking.Status() {
	case photosessionmodel.BookingStatusAccepted:
		// allowed
	case photosessionmodel.BookingStatusPendingApproval:
		return ConfirmSessionOutput{}, derrors.ErrPhotoSessionPending
	default:
		return ConfirmSessionOutput{}, derrors.ErrPhotoSessionAlreadyFinal
	}

	entry, err := s.repo.GetEntryByIDForUpdate(ctx, tx, booking.AgendaEntryID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ConfirmSessionOutput{}, derrors.Infra("agenda entry not found for booking", err)
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.confirm.get_entry_error", "agenda_entry_id", booking.AgendaEntryID(), "err", err)
		return ConfirmSessionOutput{}, derrors.Infra("failed to load agenda entry", err)
	}

	if updateErr := s.repo.UpdateBookingStatus(ctx, tx, booking.ID(), photosessionmodel.BookingStatusActive); updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return ConfirmSessionOutput{}, utils.NotFoundError("Photo session")
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("photo_session.confirm.update_booking_status_error", "booking_id", booking.ID(), "err", updateErr)
		return ConfirmSessionOutput{}, derrors.Infra("failed to update booking status", updateErr)
	}

	if updateErr := s.listingRepo.UpdateListingStatus(ctx, tx, listing.ID(), listingmodel.StatusPhotosScheduled, listingmodel.StatusPendingAvailabilityConfirm); updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return ConfirmSessionOutput{}, derrors.ErrListingNotEligible
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("photo_session.confirm.update_listing_status_error", "listing_id", listing.ID(), "err", updateErr)
		return ConfirmSessionOutput{}, derrors.Infra("failed to update listing status", updateErr)
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("photo_session.confirm.commit_error", "listing_id", listing.ID(), "err", commitErr)
		return ConfirmSessionOutput{}, derrors.Infra("failed to commit confirmation", commitErr)
	}
	committed = true

	timezone := entry.Timezone()
	loc, locErr := resolveLocation(timezone)
	if locErr != nil {
		loc = time.Local
	}
	start := booking.StartsAt().In(loc)
	end := booking.EndsAt().In(loc)

	logger.Info("photo_session.confirm.success", "booking_id", booking.ID(), "listing_id", listing.ID())

	return ConfirmSessionOutput{
		PhotoSessionID: booking.ID(),
		SlotStart:      start,
		SlotEnd:        end,
		PhotographerID: booking.PhotographerUserID(),
		ListingID:      listing.ID(),
		Status:         photosessionmodel.BookingStatusActive,
	}, nil
}

// CancelPhotoSession releases a previously reserved or confirmed session.
func (s *photoSessionService) CancelPhotoSession(ctx context.Context, input CancelSessionInput) (CancelSessionOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return CancelSessionOutput{}, derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.UserID <= 0 {
		return CancelSessionOutput{}, derrors.Auth("unauthorized")
	}
	if input.PhotoSessionID == 0 {
		return CancelSessionOutput{}, derrors.Validation("photoSessionId must be greater than zero", map[string]any{"photoSessionId": "greater_than_zero"})
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("photo_session.cancel.tx_start_error", "err", txErr)
		return CancelSessionOutput{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("photo_session.cancel.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	booking, err := s.repo.GetBookingByIDForUpdate(ctx, tx, input.PhotoSessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CancelSessionOutput{}, utils.NotFoundError("Photo session")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.cancel.get_booking_error", "photo_session_id", input.PhotoSessionID, "err", err)
		return CancelSessionOutput{}, derrors.Infra("failed to load booking", err)
	}

	listing, err := s.listingRepo.GetListingByID(ctx, tx, booking.ListingID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CancelSessionOutput{}, utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.cancel.get_listing_error", "listing_id", booking.ListingID(), "err", err)
		return CancelSessionOutput{}, derrors.Infra("failed to load listing", err)
	}

	if listing.Deleted() {
		return CancelSessionOutput{}, utils.BadRequest("listing is not available")
	}
	if listing.UserID() != input.UserID {
		return CancelSessionOutput{}, derrors.Auth("listing does not belong to user")
	}

	entry, err := s.repo.GetEntryByIDForUpdate(ctx, tx, booking.AgendaEntryID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CancelSessionOutput{}, utils.NotFoundError("Photographer agenda entry")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.cancel.get_entry_error", "agenda_entry_id", booking.AgendaEntryID(), "err", err)
		return CancelSessionOutput{}, derrors.Infra("failed to load agenda entry", err)
	}

	loc, locErr := resolveLocation(entry.Timezone())
	if locErr != nil {
		loc = time.Local
	}
	slotStart := booking.StartsAt().In(loc)
	slotEnd := booking.EndsAt().In(loc)

	var expectedStatus listingmodel.ListingStatus
	switch booking.Status() {
	case photosessionmodel.BookingStatusPendingApproval,
		photosessionmodel.BookingStatusAccepted,
		photosessionmodel.BookingStatusActive:
		if booking.Status() == photosessionmodel.BookingStatusActive {
			expectedStatus = listingmodel.StatusPhotosScheduled
		} else {
			expectedStatus = listingmodel.StatusPendingAvailabilityConfirm
		}
	default:
		return CancelSessionOutput{}, derrors.ErrPhotoSessionNotCancelable
	}

	if err := s.repo.UpdateBookingStatus(ctx, tx, booking.ID(), photosessionmodel.BookingStatusCancelled); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CancelSessionOutput{}, utils.NotFoundError("Photo session")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.cancel.update_booking_status_error", "booking_id", booking.ID(), "err", err)
		return CancelSessionOutput{}, derrors.Infra("failed to update booking status", err)
	}

	if err := s.repo.DeleteEntryByID(ctx, tx, booking.AgendaEntryID()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CancelSessionOutput{}, utils.NotFoundError("Photographer agenda entry")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.cancel.delete_entry_error", "agenda_entry_id", booking.AgendaEntryID(), "err", err)
		return CancelSessionOutput{}, derrors.Infra("failed to delete agenda entry", err)
	}

	if updateErr := s.listingRepo.UpdateListingStatus(ctx, tx, listing.ID(), listingmodel.StatusPendingPhotoScheduling, expectedStatus); updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return CancelSessionOutput{}, derrors.ErrListingNotEligible
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("photo_session.cancel.update_listing_status_error", "listing_id", listing.ID(), "err", updateErr)
		return CancelSessionOutput{}, derrors.Infra("failed to update listing status", updateErr)
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("photo_session.cancel.commit_error", "listing_id", listing.ID(), "err", commitErr)
		return CancelSessionOutput{}, derrors.Infra("failed to commit cancellation", commitErr)
	}
	committed = true

	logger.Info("photo_session.cancel.success", "booking_id", booking.ID(), "listing_id", listing.ID())

	return CancelSessionOutput{
		PhotoSessionID: booking.ID(),
		SlotStart:      slotStart,
		SlotEnd:        slotEnd,
		PhotographerID: booking.PhotographerUserID(),
		ListingID:      listing.ID(),
		ListingCode:    listing.Code(),
	}, nil
}

func (s *photoSessionService) buildAgendaSlot(entry photosessionmodel.AgendaEntryInterface, loc *time.Location) AgendaSlot {
	start := entry.StartsAt().In(loc)
	end := entry.EndsAt().In(loc)

	slot := AgendaSlot{
		EntryID:        entry.ID(),
		PhotographerID: entry.PhotographerUserID(),
		Start:          start,
		End:            end,
		Status:         photosessionmodel.SlotStatusBlocked,
		GroupID:        buildAgendaGroupID(entry, loc),
		Source:         entry.Source(),
		IsHoliday:      entry.EntryType() == photosessionmodel.AgendaEntryTypeHoliday,
		IsTimeOff:      entry.EntryType() == photosessionmodel.AgendaEntryTypeTimeOff,
		Timezone:       loc.String(),
		EntryType:      entry.EntryType(),
	}

	if !entry.Blocking() {
		slot.Status = photosessionmodel.SlotStatusAvailable
	}

	switch entry.EntryType() {
	case photosessionmodel.AgendaEntryTypePhotoSession:
		slot.Status = photosessionmodel.SlotStatusBooked
		if sourceID, ok := entry.SourceID(); ok && sourceID != nil {
			slot.SourceID = *sourceID
		}
	case photosessionmodel.AgendaEntryTypeHoliday:
		if sourceID, ok := entry.SourceID(); ok && sourceID != nil {
			slot.HolidayCalendarIDs = []uint64{*sourceID}
		}
		if reason, ok := entry.Reason(); ok {
			slot.HolidayLabels = []string{reason}
		}
	case photosessionmodel.AgendaEntryTypeTimeOff:
		if reason, ok := entry.Reason(); ok {
			slot.Reason = reason
		}
	}

	return slot
}

func (s *photoSessionService) buildDefaultBlockEntries(photographerID uint64, loc *time.Location, windowStart, windowEnd time.Time, workdayStart, workdayEnd int) []photosessionmodel.AgendaEntryInterface {
	entries := make([]photosessionmodel.AgendaEntryInterface, 0)
	for day := windowStart; day.Before(windowEnd); day = day.AddDate(0, 0, 1) {
		dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, loc)
		dayEnd := dayStart.Add(24 * time.Hour)

		// Full weekend block.
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
	if strings.TrimSpace(reason) != "" {
		entry.SetReason(strings.TrimSpace(reason))
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
				utils.LoggerFromContext(ctx).Error("photo_session.holiday.list_error", "calendar_id", calendarID, "err", err)
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

func buildAgendaGroupID(entry photosessionmodel.AgendaEntryInterface, loc *time.Location) string {
	dayKey := entry.StartsAt().In(loc).Format("2006-01-02")
	return fmt.Sprintf("%s-%s", strings.ToLower(string(entry.Source())), dayKey)
}

// UpdateSessionStatusInput contains data required to mutate a booking status.
type UpdateSessionStatusInput struct {
	SessionID      uint64
	PhotographerID uint64
	Status         string
}

// EnsureAgendaInput controls agenda bootstrap parameters.
type EnsureAgendaInput struct {
	PhotographerID    uint64
	Timezone          string
	HolidayCalendarID *uint64
	HorizonMonths     int
	WorkdayStartHour  int
	WorkdayEndHour    int
}

// TimeOffInput represents the payload to create a time-off entry.
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

// DeleteTimeOffInput represents the payload to remove a time-off entry.
type DeleteTimeOffInput struct {
	TimeOffID      uint64
	PhotographerID uint64
	Timezone       string
}

// UpdateTimeOffInput represents the payload to update a time-off entry.
type UpdateTimeOffInput struct {
	TimeOffID         uint64
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

// ListTimeOffInput captures filters for time-off listing.
type ListTimeOffInput struct {
	PhotographerID uint64
	RangeFrom      time.Time
	RangeTo        time.Time
	Page           int
	Size           int
	Timezone       string
}

// TimeOffDetailInput carries identifiers to fetch a time-off entry.
type TimeOffDetailInput struct {
	TimeOffID      uint64
	PhotographerID uint64
	Timezone       string
}

// ListTimeOffOutput aggregates paginated time-off entries.
type ListTimeOffOutput struct {
	TimeOffs []photosessionmodel.AgendaEntryInterface
	Total    int64
	Page     int
	Size     int
	Timezone string
}

// TimeOffDetailResult represents a single time-off entry alongside timezone metadata.
type TimeOffDetailResult struct {
	TimeOff  photosessionmodel.AgendaEntryInterface
	Timezone string
}

// ListAgendaInput defines the input for listing agenda entries.
type ListAgendaInput struct {
	PhotographerID     uint64
	StartDate          time.Time
	EndDate            time.Time
	Page               int
	Size               int
	Timezone           string
	HolidayCalendarIDs []uint64
}

// ListAgendaOutput describes the agenda listing result.
type ListAgendaOutput struct {
	Slots    []AgendaSlot `json:"slots"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	Size     int          `json:"size"`
	Timezone string       `json:"timezone"`
}

// AgendaSlot represents an agenda entry rendered to clients.
type AgendaSlot struct {
	EntryID            uint64                              `json:"entryId"`
	PhotographerID     uint64                              `json:"photographerId"`
	EntryType          photosessionmodel.AgendaEntryType   `json:"entryType"`
	Source             photosessionmodel.AgendaEntrySource `json:"source"`
	SourceID           uint64                              `json:"sourceId,omitempty"`
	Start              time.Time                           `json:"start"`
	End                time.Time                           `json:"end"`
	Status             photosessionmodel.SlotStatus        `json:"status"`
	GroupID            string                              `json:"groupId"`
	IsHoliday          bool                                `json:"isHoliday"`
	IsTimeOff          bool                                `json:"isTimeOff"`
	HolidayLabels      []string                            `json:"holidayLabels,omitempty"`
	HolidayCalendarIDs []uint64                            `json:"holidayCalendarIds,omitempty"`
	Reason             string                              `json:"reason,omitempty"`
	Timezone           string                              `json:"timezone"`
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

func listingAllowsPhotoSession(status listingmodel.ListingStatus) bool {
	switch status {
	case listingmodel.StatusPendingPhotoScheduling,
		listingmodel.StatusPendingAvailabilityConfirm,
		listingmodel.StatusPhotosScheduled:
		return true
	default:
		return false
	}
}

// ListAvailabilityInput encapsulates range and pagination data for availability listing.
type ListAvailabilityInput struct {
	From             *time.Time
	To               *time.Time
	Page             int
	Size             int
	Sort             string
	Period           *photosessionmodel.SlotPeriod
	Timezone         string
	WorkdayStartHour int
	WorkdayEndHour   int
}

// ListAvailabilityOutput aggregates computed availability slots.
type ListAvailabilityOutput struct {
	Slots    []AvailabilitySlot
	Total    int64
	Page     int
	Size     int
	Timezone string
}

// AvailabilitySlot represents a free window available for booking.
type AvailabilitySlot struct {
	SlotID         uint64
	PhotographerID uint64
	Start          time.Time
	End            time.Time
	Period         photosessionmodel.SlotPeriod
	SourceTimezone string
}

// ReserveSessionInput captures the necessary identifiers to reserve a photo session window.
type ReserveSessionInput struct {
	ListingID int64
	SlotID    uint64
	UserID    int64
}

// ReserveSessionOutput returns metadata about the reserved session.
type ReserveSessionOutput struct {
	PhotoSessionID uint64
	SlotID         uint64
	SlotStart      time.Time
	SlotEnd        time.Time
	PhotographerID uint64
	ListingID      int64
}

// ConfirmSessionInput holds data required to confirm a reserved session.
type ConfirmSessionInput struct {
	ListingID      int64
	PhotoSessionID uint64
	UserID         int64
}

// ConfirmSessionOutput reports the confirmed session metadata.
type ConfirmSessionOutput struct {
	PhotoSessionID uint64
	SlotStart      time.Time
	SlotEnd        time.Time
	PhotographerID uint64
	ListingID      int64
	Status         photosessionmodel.BookingStatus
}

// CancelSessionInput captures identifiers needed to cancel an existing session.
type CancelSessionInput struct {
	PhotoSessionID uint64
	UserID         int64
}

// CancelSessionOutput reports metadata about a cancelled session.
type CancelSessionOutput struct {
	PhotoSessionID uint64
	SlotStart      time.Time
	SlotEnd        time.Time
	PhotographerID uint64
	ListingID      int64
	ListingCode    uint32
}

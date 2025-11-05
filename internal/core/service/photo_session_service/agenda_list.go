package photosessionservices

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListAgenda returns agenda entries within the interval.
func (s *photoSessionService) ListAgenda(ctx context.Context, input ListAgendaInput) (ListAgendaOutput, error) {
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "service.ListAgenda")
	if err != nil {
		return ListAgendaOutput{}, derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err := validateListAgendaInput(input); err != nil {
		return ListAgendaOutput{}, err
	}

	loc := input.Location
	if loc == nil {
		loc = time.UTC
	}

	startLocal := utils.ConvertToLocation(input.StartDate, loc)
	endLocal := utils.ConvertToLocation(input.EndDate, loc)

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
		logger.Error("service.list_agenda.tx_start_error", "err", err)
		return ListAgendaOutput{}, derrors.Wrap(err, derrors.KindInfra, "failed to start transaction")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("service.list_agenda.tx_rollback_error", "err", rbErr)
		}
	}()

	profile, profileErr := s.loadPhotographerLocation(ctx, tx, input.PhotographerID)
	if profileErr != nil {
		return ListAgendaOutput{}, profileErr
	}

	entries, err := s.repo.ListEntriesByRange(ctx, tx, input.PhotographerID, utils.ConvertToUTC(startLocal), utils.ConvertToUTC(endLocal), input.EntryType)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.list_agenda.repo_error", "photographer_id", input.PhotographerID, "err", err)
		return ListAgendaOutput{}, derrors.Wrap(err, derrors.KindInfra, "failed to list agenda entries")
	}

	// Carregar bookings associados para popular photoSessionId
	bookingsMap, err := s.loadBookingsForEntries(ctx, tx, entries)
	if err != nil {
		return ListAgendaOutput{}, derrors.Wrap(err, derrors.KindInfra, "failed to load bookings for agenda entries")
	}

	occupied := make(map[string]struct{}, len(entries))
	slots := make([]AgendaSlot, 0, len(entries))
	for _, entry := range entries {
		startLocal := entry.StartsAt().In(loc)
		endLocal := entry.EndsAt().In(loc)
		occupied[agendaSlotKey(entry.EntryType(), entry.Source(), startLocal, endLocal)] = struct{}{}
		slots = append(slots, s.buildAgendaSlot(entry, loc, bookingsMap))
	}

	// Only add holiday slots if no filter or filter includes HOLIDAY
	if input.EntryType == nil || *input.EntryType == photosessionmodel.AgendaEntryTypeHoliday {
		holidaySlots, holidayErr := s.fetchHolidaySlots(ctx, input.PhotographerID, loc, profile, startLocal, endLocal, occupied)
		if holidayErr != nil {
			return ListAgendaOutput{}, holidayErr
		}
		slots = append(slots, holidaySlots...)
	}

	// Only add non-working slots (blocks) if no filter or filter includes BLOCK
	if input.EntryType == nil || *input.EntryType == photosessionmodel.AgendaEntryTypeBlock {
		nonWorkingSlots := s.buildNonWorkingSlots(input.PhotographerID, loc, startLocal, endLocal, occupied)
		slots = append(slots, nonWorkingSlots...)
	}

	sort.Slice(slots, func(i, j int) bool {
		if slots[i].Start.Equal(slots[j].Start) {
			if slots[i].End.Equal(slots[j].End) {
				return slots[i].EntryID < slots[j].EntryID
			}
			return slots[i].End.Before(slots[j].End)
		}
		return slots[i].Start.Before(slots[j].Start)
	})

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

func (s *photoSessionService) buildAgendaSlot(entry photosessionmodel.AgendaEntryInterface, loc *time.Location, bookingsMap map[uint64]uint64) AgendaSlot {
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

		// Popular photoSessionId se disponível no map
		if bookingID, found := bookingsMap[entry.ID()]; found {
			slot.PhotoSessionID = &bookingID
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

func buildAgendaGroupID(entry photosessionmodel.AgendaEntryInterface, loc *time.Location) string {
	dayKey := entry.StartsAt().In(loc).Format("2006-01-02")
	return fmt.Sprintf("%s-%s", strings.ToLower(string(entry.Source())), dayKey)
}

// loadBookingsForEntries retrieves bookings associated with agenda entries.
// Returns a map of agendaEntryID -> bookingID for quick lookup.
func (s *photoSessionService) loadBookingsForEntries(ctx context.Context, tx *sql.Tx, entries []photosessionmodel.AgendaEntryInterface) (map[uint64]uint64, error) {
	logger := utils.LoggerFromContext(ctx)

	// Coleta IDs de entradas do tipo PHOTO_SESSION com source BOOKING
	entryIDs := make([]uint64, 0)
	for _, entry := range entries {
		if entry.EntryType() == photosessionmodel.AgendaEntryTypePhotoSession &&
			entry.Source() == photosessionmodel.AgendaEntrySourceBooking {
			entryIDs = append(entryIDs, entry.ID())
		}
	}

	if len(entryIDs) == 0 {
		return make(map[uint64]uint64), nil
	}

	// Buscar bookings associados
	bookingsMap := make(map[uint64]uint64, len(entryIDs))
	for _, entryID := range entryIDs {
		booking, err := s.repo.FindBookingByAgendaEntry(ctx, tx, entryID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Booking não encontrado - log de warning, não bloqueia
				logger.Warn("agenda.load_bookings.booking_not_found", "agenda_entry_id", entryID)
				continue
			}
			utils.SetSpanError(ctx, err)
			logger.Error("agenda.load_bookings.repo_error", "agenda_entry_id", entryID, "err", err)
			return nil, fmt.Errorf("failed to load booking for entry %d: %w", entryID, err)
		}
		bookingsMap[entryID] = booking.ID()
	}

	return bookingsMap, nil
}

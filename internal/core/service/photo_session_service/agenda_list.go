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

	sortField := sanitizeAgendaSortField(input.SortField)
	sortOrder := sanitizeAgendaSortOrder(input.SortOrder)

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
		return compareAgendaSlots(slots[i], slots[j], sortField, sortOrder)
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

// Internal struct to hold booking and listing details
type bookingDetails struct {
	BookingID uint64
	Listing   *ListingInfo
	Status    photosessionmodel.BookingStatus
}

func sanitizeAgendaSortField(field AgendaSortField) AgendaSortField {
	switch field {
	case AgendaSortFieldEndDate, AgendaSortFieldEntryType, AgendaSortFieldStartDate:
		return field
	default:
		return defaultAgendaSortField
	}
}

func sanitizeAgendaSortOrder(order AgendaSortOrder) AgendaSortOrder {
	switch order {
	case AgendaSortOrderDesc:
		return AgendaSortOrderDesc
	default:
		return defaultAgendaSortOrder
	}
}

func compareAgendaSlots(a, b AgendaSlot, sortField AgendaSortField, sortOrder AgendaSortOrder) bool {
	if sortOrder == AgendaSortOrderDesc {
		return compareAgendaSlotsAsc(b, a, sortField)
	}
	return compareAgendaSlotsAsc(a, b, sortField)
}

func compareAgendaSlotsAsc(a, b AgendaSlot, sortField AgendaSortField) bool {
	switch sortField {
	case AgendaSortFieldEndDate:
		return compareTimeThenID(a.End, b.End, a.EntryID, b.EntryID)
	case AgendaSortFieldEntryType:
		if a.EntryType == b.EntryType {
			return compareTimeThenID(a.Start, b.Start, a.EntryID, b.EntryID)
		}
		return string(a.EntryType) < string(b.EntryType)
	default:
		return compareTimeThenID(a.Start, b.Start, a.EntryID, b.EntryID)
	}
}

func compareTimeThenID(aTime, bTime time.Time, aID, bID uint64) bool {
	if aTime.Equal(bTime) {
		return aID < bID
	}
	return aTime.Before(bTime)
}

func (s *photoSessionService) buildAgendaSlot(entry photosessionmodel.AgendaEntryInterface, loc *time.Location, bookingsMap map[uint64]bookingDetails) AgendaSlot {
	start := entry.StartsAt().In(loc)
	end := entry.EndsAt().In(loc)

	// Default status based on blocking nature
	defaultStatus := photosessionmodel.SlotStatusBlocked
	if !entry.Blocking() {
		defaultStatus = photosessionmodel.SlotStatusAvailable
	}

	slot := AgendaSlot{
		EntryID:        entry.ID(),
		PhotographerID: entry.PhotographerUserID(),
		Start:          start,
		End:            end,
		Status:         string(defaultStatus),
		GroupID:        buildAgendaGroupID(entry, loc),
		Source:         entry.Source(),
		IsHoliday:      entry.EntryType() == photosessionmodel.AgendaEntryTypeHoliday,
		IsTimeOff:      entry.EntryType() == photosessionmodel.AgendaEntryTypeTimeOff,
		Timezone:       loc.String(),
		EntryType:      entry.EntryType(),
	}

	switch entry.EntryType() {
	case photosessionmodel.AgendaEntryTypePhotoSession:
		// Default to BOOKED if no specific booking status is found
		slot.Status = string(photosessionmodel.SlotStatusBooked)

		if sourceID, ok := entry.SourceID(); ok && sourceID != nil {
			slot.SourceID = *sourceID
		}

		// Popular photoSessionId e Listing se disponível no map
		if details, found := bookingsMap[entry.ID()]; found {
			slot.PhotoSessionID = &details.BookingID
			slot.Listing = details.Listing
			// Sobrescreve com o status real do booking (ex: DONE, ACCEPTED)
			slot.Status = string(details.Status)
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
// Returns a map of agendaEntryID -> bookingDetails for quick lookup.
func (s *photoSessionService) loadBookingsForEntries(ctx context.Context, tx *sql.Tx, entries []photosessionmodel.AgendaEntryInterface) (map[uint64]bookingDetails, error) {
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
		return make(map[uint64]bookingDetails), nil
	}

	// Buscar bookings associados
	bookingsMap := make(map[uint64]bookingDetails, len(entryIDs))
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

		// Fetch listing details
		listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, booking.ListingIdentityID())
		if err != nil {
			// Log warning and continue, or return error depending on strictness.
			// Suggest logging warning to not break agenda if listing is missing.
			logger.Warn("agenda.load_bookings.listing_not_found", "listing_id", booking.ListingIdentityID(), "err", err)
			bookingsMap[entryID] = bookingDetails{
				BookingID: booking.ID(),
				Status:    booking.Status(),
			}
			continue
		}

		bookingsMap[entryID] = bookingDetails{
			BookingID: booking.ID(),
			Status:    booking.Status(),
			Listing: &ListingInfo{
				ListingIdentityID: listing.IdentityID(),
				Code:              listing.Code(),
				Title:             listing.Title(),
				ZipCode:           listing.ZipCode(),
				Street:            listing.Street(),
				Number:            listing.Number(),
				Complement:        listing.Complement(),
				Neighborhood:      listing.Neighborhood(),
				City:              listing.City(),
				State:             listing.State(),
				Status:            listing.Status().String(),
			},
		}
	}

	return bookingsMap, nil
}

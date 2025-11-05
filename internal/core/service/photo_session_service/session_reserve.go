package photosessionservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

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

	if updateErr := s.listingRepo.UpdateListingStatus(ctx, tx, listing.ID(), listingmodel.StatusPendingPhotoConfirmation, listing.Status()); updateErr != nil {
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

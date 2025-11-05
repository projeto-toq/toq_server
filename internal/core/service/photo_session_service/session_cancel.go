package photosessionservices

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

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
			expectedStatus = listingmodel.StatusPendingPhotoConfirmation
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

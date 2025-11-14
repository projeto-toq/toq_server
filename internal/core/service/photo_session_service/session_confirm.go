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

	listing, err := s.listingRepo.GetListingVersionByID(ctx, tx, input.ListingID)
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

	if updateErr := s.listingRepo.UpdateListingStatus(ctx, tx, listing.ID(), listingmodel.StatusPhotosScheduled, listingmodel.StatusPendingPhotoConfirmation); updateErr != nil {
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

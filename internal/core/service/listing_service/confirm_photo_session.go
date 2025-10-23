package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"time"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) ConfirmPhotoSession(ctx context.Context, input ConfirmPhotoSessionInput) (output ConfirmPhotoSessionOutput, err error) {
	if input.ListingID <= 0 || input.SlotID == 0 || input.ReservationToken == "" {
		return output, utils.BadRequest("listingId, slotId and reservationToken are required")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return output, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	userID, userErr := ls.gsi.GetUserIDFromContext(ctx)
	if userErr != nil {
		return output, userErr
	}

	tx, txErr := ls.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.photo_session.confirm.start_tx_error", "err", txErr)
		return output, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.photo_session.confirm.rollback_error", "err", rbErr)
			}
		}
	}()

	listing, repoErr := ls.listingRepository.GetListingByID(ctx, tx, input.ListingID)
	if repoErr != nil {
		if errors.Is(repoErr, sql.ErrNoRows) {
			return output, utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, repoErr)
		logger.Error("listing.photo_session.confirm.get_listing_error", "err", repoErr, "listing_id", input.ListingID)
		return output, utils.InternalError("")
	}

	if listing.Deleted() {
		return output, utils.BadRequest("listing is not available")
	}

	if listing.UserID() != userID {
		return output, utils.AuthorizationError("listing does not belong to current user")
	}

	if !listingEligibleForPhotoSession(listing.Status()) {
		return output, derrors.ErrListingNotEligible
	}

	slot, slotErr := ls.photoSessionRepo.GetSlotForUpdate(ctx, tx, input.SlotID)
	if slotErr != nil {
		if errors.Is(slotErr, sql.ErrNoRows) {
			return output, derrors.ErrSlotUnavailable
		}
		utils.SetSpanError(ctx, slotErr)
		logger.Error("listing.photo_session.confirm.get_slot_error", "err", slotErr, "slot_id", input.SlotID)
		return output, utils.InternalError("")
	}

	if slot.Status() != photosessionmodel.SlotStatusReserved {
		return output, derrors.ErrSlotUnavailable
	}

	token := slot.ReservationToken()
	if token == nil || *token != input.ReservationToken {
		return output, derrors.ErrSlotUnavailable
	}

	expiresAt := slot.ReservedUntil()
	if expiresAt == nil || time.Now().UTC().After(*expiresAt) {
		return output, derrors.ErrReservationExpired
	}

	booking := photosessionmodel.NewPhotoSessionBooking()
	booking.SetSlotID(slot.ID())
	booking.SetListingID(input.ListingID)
	booking.SetScheduledStart(slot.SlotStart())
	booking.SetScheduledEnd(slot.SlotEnd())
	booking.SetStatus(photosessionmodel.BookingStatusActive)

	bookingID, insertErr := ls.photoSessionRepo.InsertBooking(ctx, tx, booking)
	if insertErr != nil {
		utils.SetSpanError(ctx, insertErr)
		logger.Error("listing.photo_session.confirm.insert_booking_error", "err", insertErr, "slot_id", input.SlotID)
		return output, utils.InternalError("")
	}

	if markErr := ls.photoSessionRepo.MarkSlotBooked(ctx, tx, input.SlotID, time.Now().UTC()); markErr != nil {
		if errors.Is(markErr, sql.ErrNoRows) {
			return output, derrors.ErrSlotUnavailable
		}
		utils.SetSpanError(ctx, markErr)
		logger.Error("listing.photo_session.confirm.update_slot_error", "err", markErr, "slot_id", input.SlotID)
		return output, utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.photo_session.confirm.commit_error", "err", cmErr)
		return output, utils.InternalError("")
	}

	logger.Info("listing.photo_session.confirm.success", "booking_id", bookingID, "slot_id", input.SlotID, "listing_id", input.ListingID, "user_id", userID)

	return ConfirmPhotoSessionOutput{
		PhotoSessionID: bookingID,
		SlotID:         input.SlotID,
		ScheduledStart: slot.SlotStart(),
		ScheduledEnd:   slot.SlotEnd(),
	}, nil
}

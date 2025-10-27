package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) CancelPhotoSession(ctx context.Context, input CancelPhotoSessionInput) (err error) {
	if input.PhotoSessionID == 0 {
		return utils.BadRequest("photoSessionId is required")
	}

	ctx, spanEnd, genErr := utils.GenerateTracer(ctx)
	if genErr != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	userID, userErr := ls.gsi.GetUserIDFromContext(ctx)
	if userErr != nil {
		return userErr
	}

	tx, txErr := ls.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.photo_session.cancel.start_tx_error", "err", txErr)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.photo_session.cancel.rollback_error", "err", rbErr)
			}
		}
	}()

	booking, repoErr := ls.photoSessionRepo.GetBookingForUpdate(ctx, tx, input.PhotoSessionID)
	if repoErr != nil {
		if errors.Is(repoErr, sql.ErrNoRows) {
			return utils.NotFoundError("Photo session")
		}
		utils.SetSpanError(ctx, repoErr)
		logger.Error("listing.photo_session.cancel.get_booking_error", "err", repoErr, "photo_session_id", input.PhotoSessionID)
		return utils.InternalError("")
	}

	switch booking.Status() {
	case photosessionmodel.BookingStatusPendingApproval,
		photosessionmodel.BookingStatusAccepted,
		photosessionmodel.BookingStatusActive:
		// allowed
	default:
		return derrors.ErrPhotoSessionNotCancelable
	}

	listing, listingErr := ls.listingRepository.GetListingByID(ctx, tx, booking.ListingID())
	if listingErr != nil {
		if errors.Is(listingErr, sql.ErrNoRows) {
			return utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, listingErr)
		logger.Error("listing.photo_session.cancel.get_listing_error", "err", listingErr, "listing_id", booking.ListingID())
		return utils.InternalError("")
	}

	if listing.UserID() != userID {
		return utils.AuthorizationError("listing does not belong to current user")
	}

	slot, slotErr := ls.photoSessionRepo.GetSlotForUpdate(ctx, tx, booking.SlotID())
	if slotErr != nil {
		if errors.Is(slotErr, sql.ErrNoRows) {
			return utils.NotFoundError("Photographer slot")
		}
		utils.SetSpanError(ctx, slotErr)
		logger.Error("listing.photo_session.cancel.get_slot_error", "err", slotErr, "slot_id", booking.SlotID())
		return utils.InternalError("")
	}

	photographerPhone := ""
	if slot.PhotographerUserID() != 0 && ls.userRepository != nil {
		if photographer, photographerErr := ls.userRepository.GetUserByID(ctx, tx, int64(slot.PhotographerUserID())); photographerErr != nil {
			if !errors.Is(photographerErr, sql.ErrNoRows) {
				utils.SetSpanError(ctx, photographerErr)
				logger.Error("listing.photo_session.cancel.get_photographer_error", "err", photographerErr, "photographer_id", slot.PhotographerUserID())
				return utils.InternalError("")
			}
		} else if phone := strings.TrimSpace(photographer.GetPhoneNumber()); phone != "" {
			photographerPhone = phone
		}
	}

	if updateErr := ls.photoSessionRepo.UpdateBookingStatus(ctx, tx, input.PhotoSessionID, photosessionmodel.BookingStatusCancelled); updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return utils.NotFoundError("Photo session")
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("listing.photo_session.cancel.update_booking_error", "err", updateErr, "photo_session_id", input.PhotoSessionID)
		return utils.InternalError("")
	}

	if markErr := ls.photoSessionRepo.MarkSlotAvailable(ctx, tx, booking.SlotID()); markErr != nil {
		if errors.Is(markErr, sql.ErrNoRows) {
			return utils.NotFoundError("Photographer slot")
		}
		utils.SetSpanError(ctx, markErr)
		logger.Error("listing.photo_session.cancel.update_slot_error", "err", markErr, "slot_id", booking.SlotID())
		return utils.InternalError("")
	}

	var expectedStatus listingmodel.ListingStatus
	switch booking.Status() {
	case photosessionmodel.BookingStatusActive:
		expectedStatus = listingmodel.StatusPhotosScheduled
	default:
		expectedStatus = listingmodel.StatusPendingAvailabilityConfirm
	}

	if updateErr := ls.listingRepository.UpdateListingStatus(ctx, tx, listing.ID(), listingmodel.StatusPendingPhotoScheduling, expectedStatus); updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return derrors.ErrListingNotEligible
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("listing.photo_session.cancel.update_listing_status_error", "err", updateErr, "listing_id", listing.ID())
		return utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.photo_session.cancel.commit_error", "err", cmErr)
		return utils.InternalError("")
	}

	logger.Info("listing.photo_session.cancel.success", "photo_session_id", input.PhotoSessionID, "listing_id", booking.ListingID(), "user_id", userID)

	ls.sendPhotographerCancellationSMS(ctx, photographerPhone, slot.SlotStart(), listing.Code())

	return nil
}

package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) ReservePhotoSession(ctx context.Context, input ReservePhotoSessionInput) (output ReservePhotoSessionOutput, err error) {
	if input.ListingID <= 0 || input.SlotID == 0 {
		return output, utils.BadRequest("listingId and slotId are required")
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
		logger.Error("listing.photo_session.reserve.start_tx_error", "err", txErr)
		return output, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.photo_session.reserve.rollback_error", "err", rbErr)
			}
		}
	}()

	listing, repoErr := ls.listingRepository.GetListingByID(ctx, tx, input.ListingID)
	if repoErr != nil {
		if errors.Is(repoErr, sql.ErrNoRows) {
			return output, utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, repoErr)
		logger.Error("listing.photo_session.reserve.get_listing_error", "err", repoErr, "listing_id", input.ListingID)
		return output, utils.InternalError("")
	}

	if listing.Deleted() {
		return output, utils.BadRequest("listing is not available")
	}

	if listing.UserID() != userID {
		return output, utils.AuthorizationError("listing does not belong to current user")
	}

	if listing.Status() == listingmodel.StatusPendingAvailabilityConfirm {
		return output, derrors.ErrListingNotEligible
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
		logger.Error("listing.photo_session.reserve.get_slot_error", "err", slotErr, "slot_id", input.SlotID)
		return output, utils.InternalError("")
	}

	if slot.Status() != photosessionmodel.SlotStatusAvailable {
		return output, derrors.ErrSlotUnavailable
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	if slot.SlotDate().Before(today) {
		return output, derrors.ErrSlotUnavailable
	}

	reservationToken := uuid.New().String()
	expiresAt := time.Now().UTC().Add(reservationHoldTTL)

	photographerPhone := ""
	if slot.PhotographerUserID() != 0 && ls.userRepository != nil {
		if photographer, photographerErr := ls.userRepository.GetUserByID(ctx, tx, int64(slot.PhotographerUserID())); photographerErr != nil {
			if !errors.Is(photographerErr, sql.ErrNoRows) {
				utils.SetSpanError(ctx, photographerErr)
				logger.Error("listing.photo_session.reserve.get_photographer_error", "err", photographerErr, "photographer_id", slot.PhotographerUserID())
				return output, utils.InternalError("")
			}
		} else if phone := strings.TrimSpace(photographer.GetPhoneNumber()); phone != "" {
			photographerPhone = phone
		}
	}

	if markErr := ls.photoSessionRepo.MarkSlotReserved(ctx, tx, input.SlotID, reservationToken, expiresAt); markErr != nil {
		if errors.Is(markErr, sql.ErrNoRows) {
			return output, derrors.ErrSlotUnavailable
		}
		utils.SetSpanError(ctx, markErr)
		logger.Error("listing.photo_session.reserve.update_slot_error", "err", markErr, "slot_id", input.SlotID)
		return output, utils.InternalError("")
	}

	booking := photosessionmodel.NewPhotoSessionBooking()
	booking.SetSlotID(slot.ID())
	booking.SetListingID(input.ListingID)
	booking.SetScheduledStart(slot.SlotStart())
	booking.SetScheduledEnd(slot.SlotEnd())
	booking.SetStatus(photosessionmodel.BookingStatusPendingApproval)

	bookingID, insertErr := ls.photoSessionRepo.InsertBooking(ctx, tx, booking)
	if insertErr != nil {
		utils.SetSpanError(ctx, insertErr)
		logger.Error("listing.photo_session.reserve.insert_booking_error", "err", insertErr, "slot_id", input.SlotID)
		return output, utils.InternalError("")
	}

	if updateErr := ls.listingRepository.UpdateListingStatus(ctx, tx, input.ListingID, listingmodel.StatusPendingAvailabilityConfirm, listing.Status()); updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			utils.SetSpanError(ctx, updateErr)
			logger.Warn("listing.photo_session.reserve.status_conflict", "err", updateErr, "listing_id", input.ListingID, "current_status", listing.Status())
			return output, derrors.ErrListingNotEligible
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("listing.photo_session.reserve.update_listing_status_error", "err", updateErr, "listing_id", input.ListingID)
		return output, utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.photo_session.reserve.commit_error", "err", cmErr)
		return output, utils.InternalError("")
	}

	logger.Info("listing.photo_session.reserve.success", "listing_id", input.ListingID, "slot_id", input.SlotID, "booking_id", bookingID, "user_id", userID)

	ls.sendPhotographerReservationSMS(ctx, photographerPhone, slot.SlotStart(), slot.SlotEnd(), listing.Code())

	return ReservePhotoSessionOutput{
		SlotID:           input.SlotID,
		SlotStart:        slot.SlotStart(),
		SlotEnd:          slot.SlotEnd(),
		ReservationToken: reservationToken,
		ExpiresAt:        expiresAt,
		PhotoSessionID:   bookingID,
	}, nil
}

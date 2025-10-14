package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
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

	if markErr := ls.photoSessionRepo.MarkSlotReserved(ctx, tx, input.SlotID, reservationToken, expiresAt); markErr != nil {
		if errors.Is(markErr, sql.ErrNoRows) {
			return output, derrors.ErrSlotUnavailable
		}
		utils.SetSpanError(ctx, markErr)
		logger.Error("listing.photo_session.reserve.update_slot_error", "err", markErr, "slot_id", input.SlotID)
		return output, utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.photo_session.reserve.commit_error", "err", cmErr)
		return output, utils.InternalError("")
	}

	logger.Info("listing.photo_session.reserve.success", "listing_id", input.ListingID, "slot_id", input.SlotID, "user_id", userID)

	return ReservePhotoSessionOutput{
		SlotID:           input.SlotID,
		ReservationToken: reservationToken,
		ExpiresAt:        expiresAt,
	}, nil
}

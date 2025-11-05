package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
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

	tx, txErr := ls.gsi.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.photo_session.reserve.start_ro_tx_error", "err", txErr)
		return output, utils.InternalError("")
	}
	defer func() {
		if err != nil && tx != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.photo_session.reserve.ro_rollback_error", "err", rbErr)
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

	if listing.Status() == listingmodel.StatusPendingPhotoConfirmation {
		return output, derrors.ErrListingNotEligible
	}

	if !listingEligibleForPhotoSession(listing.Status()) {
		return output, derrors.ErrListingNotEligible
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.photo_session.reserve.ro_commit_error", "err", cmErr)
		return output, utils.InternalError("")
	}
	tx = nil

	reserveOutput, reserveErr := ls.photoSessionSvc.ReservePhotoSession(ctx, photosessionservices.ReserveSessionInput{
		ListingID: input.ListingID,
		SlotID:    input.SlotID,
		UserID:    userID,
	})
	if reserveErr != nil {
		return output, reserveErr
	}

	photographerPhone := ""
	if reserveOutput.PhotographerID != 0 && ls.userRepository != nil {
		tx, txErr = ls.gsi.StartReadOnlyTransaction(ctx)
		if txErr != nil {
			utils.SetSpanError(ctx, txErr)
			logger.Error("listing.photo_session.reserve.photographer_ro_tx_error", "err", txErr)
			return output, utils.InternalError("")
		}
		phone, fetchErr := ls.fetchPhotographerPhone(ctx, tx, reserveOutput.PhotographerID)
		if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("listing.photo_session.reserve.photographer_ro_rollback_error", "err", rbErr)
		}
		if fetchErr != nil {
			utils.SetSpanError(ctx, fetchErr)
			logger.Error("listing.photo_session.reserve.get_photographer_error", "err", fetchErr, "photographer_id", reserveOutput.PhotographerID)
			return output, utils.InternalError("")
		}
		photographerPhone = phone
	}

	logger.Info("listing.photo_session.reserve.success", "listing_id", input.ListingID, "slot_id", input.SlotID, "booking_id", reserveOutput.PhotoSessionID, "user_id", userID)

	ls.sendPhotographerReservationSMS(ctx, photographerPhone, reserveOutput.SlotStart, reserveOutput.SlotEnd, listing.Code())

	return ReservePhotoSessionOutput{
		SlotID:         reserveOutput.SlotID,
		SlotStart:      reserveOutput.SlotStart,
		SlotEnd:        reserveOutput.SlotEnd,
		PhotoSessionID: reserveOutput.PhotoSessionID,
		PhotographerID: reserveOutput.PhotographerID,
	}, nil
}

func (ls *listingService) fetchPhotographerPhone(ctx context.Context, tx *sql.Tx, photographerID uint64) (string, error) {
	if photographerID == 0 || ls.userRepository == nil {
		return "", nil
	}

	photographer, err := ls.userRepository.GetUserByID(ctx, tx, int64(photographerID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return strings.TrimSpace(photographer.GetPhoneNumber()), nil
}

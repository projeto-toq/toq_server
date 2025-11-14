package listingservices

import (
	"context"

	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
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

	cancelOutput, cancelErr := ls.photoSessionSvc.CancelPhotoSession(ctx, photosessionservices.CancelSessionInput{
		PhotoSessionID: input.PhotoSessionID,
		UserID:         userID,
	})
	if cancelErr != nil {
		return cancelErr
	}

	photographerPhone := ""
	if cancelOutput.PhotographerID != 0 && ls.userRepository != nil {
		tx, txErr := ls.gsi.StartReadOnlyTransaction(ctx)
		if txErr != nil {
			utils.SetSpanError(ctx, txErr)
			logger.Error("listing.photo_session.cancel.photographer_ro_tx_error", "err", txErr)
			return utils.InternalError("")
		}
		phone, fetchErr := ls.fetchPhotographerPhone(ctx, tx, cancelOutput.PhotographerID)
		if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("listing.photo_session.cancel.photographer_ro_rollback_error", "err", rbErr)
		}
		if fetchErr != nil {
			utils.SetSpanError(ctx, fetchErr)
			logger.Error("listing.photo_session.cancel.get_photographer_error", "err", fetchErr, "photographer_id", cancelOutput.PhotographerID)
			return utils.InternalError("")
		}
		photographerPhone = phone
	}

	logger.Info("listing.photo_session.cancel.success", "photo_session_id", cancelOutput.PhotoSessionID, "listing_identity_id", cancelOutput.ListingIdentityID, "user_id", userID)

	ls.sendPhotographerCancellationSMS(ctx, photographerPhone, cancelOutput.SlotStart, cancelOutput.ListingCode)

	return nil
}

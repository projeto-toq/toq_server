package listingservices

import (
	"context"

	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) ConfirmPhotoSession(ctx context.Context, input ConfirmPhotoSessionInput) (output ConfirmPhotoSessionOutput, err error) {
	if input.ListingID <= 0 || input.PhotoSessionID == 0 {
		return output, utils.BadRequest("listingId and photoSessionId are required")
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

	confirmOutput, confirmErr := ls.photoSessionSvc.ConfirmPhotoSession(ctx, photosessionservices.ConfirmSessionInput{
		ListingID:      input.ListingID,
		PhotoSessionID: input.PhotoSessionID,
		UserID:         userID,
	})
	if confirmErr != nil {
		return output, confirmErr
	}

	logger.Info("listing.photo_session.confirm.success", "photo_session_id", confirmOutput.PhotoSessionID, "listing_id", confirmOutput.ListingID, "user_id", userID)

	return ConfirmPhotoSessionOutput{
		PhotoSessionID: confirmOutput.PhotoSessionID,
		ScheduledStart: confirmOutput.SlotStart,
		ScheduledEnd:   confirmOutput.SlotEnd,
		Status:         string(confirmOutput.Status),
	}, nil
}

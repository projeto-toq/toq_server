package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	storagemodel "github.com/projeto-toq/toq_server/internal/core/model/storage_model"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) ReservePhotoSession(ctx context.Context, input ReservePhotoSessionInput) (output ReservePhotoSessionOutput, err error) {
	if input.ListingIdentityID <= 0 || input.SlotID == 0 {
		return output, utils.BadRequest("listingIdentityId and slotID are required")
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

	// Get listing identity to validate ownership BEFORE fetching active version
	identity, err := ls.listingRepository.GetListingIdentityByID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return output, utils.NotFoundError("listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.photo_session.reserve.get_identity_error", "err", err, "identity_id", input.ListingIdentityID)
		return output, utils.InternalError("")
	}

	// Validate ownership using identity
	if identity.UserID != userID {
		logger.Warn("unauthorized_reserve_attempt",
			"listing_identity_id", input.ListingIdentityID,
			"requester_user_id", userID,
			"owner_user_id", identity.UserID)
		return output, utils.AuthorizationError("listing does not belong to current user")
	}

	// Load active listing version to validate status
	listing, err := ls.listingRepository.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return output, utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.photo_session.reserve.get_listing_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return output, utils.InternalError("")
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
		ListingIdentityID: listing.IdentityID(),
		SlotID:            input.SlotID,
		UserID:            userID,
	})
	if reserveErr != nil {
		return output, reserveErr
	}

	photographerSummary := PhotographerSummary{}
	if reserveOutput.PhotographerID != 0 && ls.userRepository != nil {
		tx, txErr = ls.gsi.StartReadOnlyTransaction(ctx)
		if txErr != nil {
			utils.SetSpanError(ctx, txErr)
			logger.Error("listing.photo_session.reserve.photographer_ro_tx_error", "err", txErr)
			return output, utils.InternalError("")
		}
		summary, fetchErr := ls.fetchPhotographerProfile(ctx, tx, reserveOutput.PhotographerID)
		if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("listing.photo_session.reserve.photographer_ro_rollback_error", "err", rbErr)
		}
		if fetchErr != nil {
			utils.SetSpanError(ctx, fetchErr)
			logger.Error("listing.photo_session.reserve.get_photographer_error", "err", fetchErr, "photographer_id", reserveOutput.PhotographerID)
			return output, utils.InternalError("")
		}
		photographerSummary = summary
	}

	logger.Info("listing.photo_session.reserve.success", "listing_identity_id", input.ListingIdentityID, "listing_version_id", listing.ID(), "slot_id", input.SlotID, "booking_id", reserveOutput.PhotoSessionID, "user_id", userID)

	ls.sendPhotographerReservationSMS(ctx, photographerSummary.PhoneNumber, reserveOutput.SlotStart, reserveOutput.SlotEnd, listing.Code())

	return ReservePhotoSessionOutput{
		SlotID:         reserveOutput.SlotID,
		SlotStart:      reserveOutput.SlotStart,
		SlotEnd:        reserveOutput.SlotEnd,
		PhotoSessionID: reserveOutput.PhotoSessionID,
		Photographer:   photographerSummary,
	}, nil
}

func (ls *listingService) fetchPhotographerProfile(ctx context.Context, tx *sql.Tx, photographerID uint64) (PhotographerSummary, error) {
	if photographerID == 0 || ls.userRepository == nil {
		return PhotographerSummary{}, nil
	}

	photographer, err := ls.userRepository.GetUserByID(ctx, tx, int64(photographerID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PhotographerSummary{}, nil
		}
		return PhotographerSummary{}, err
	}

	phone := strings.TrimSpace(photographer.GetPhoneNumber())
	photoURL := ""
	if ls.gcs != nil {
		if url, genErr := ls.gcs.GeneratePhotoDownloadURL(int64(photographerID), storagemodel.PhotoMedium); genErr == nil {
			photoURL = url
		}
	}

	return PhotographerSummary{
		ID:          photographerID,
		FullName:    photographer.GetFullName(),
		PhoneNumber: phone,
		PhotoURL:    photoURL,
	}, nil
}

// fetchPhotographerPhone retorna apenas o telefone para cenários que não precisam do perfil completo.
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

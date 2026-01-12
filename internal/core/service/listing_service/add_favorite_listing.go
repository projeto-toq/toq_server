package listingservices

import (
	"context"
	"database/sql"
	"errors"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// AddFavoriteListing links the authenticated user to a published listing.
// Business rules:
//   - Only published listings can be favorited
//   - Owners cannot favorite their own listings
//   - Operation is idempotent
func (ls *listingService) AddFavoriteListing(ctx context.Context, listingIdentityID int64) (err error) {
	if listingIdentityID <= 0 {
		return utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}

	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	userID, err := ls.gsi.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	tx, txErr := ls.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.favorite.add.tx_start_error", "err", txErr)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.favorite.add.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	identity, idErr := ls.listingRepository.GetListingIdentityByID(ctx, tx, listingIdentityID)
	if idErr != nil {
		if errors.Is(idErr, sql.ErrNoRows) {
			return utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, idErr)
		logger.Error("listing.favorite.add.identity_error", "err", idErr, "listing_identity_id", listingIdentityID)
		return utils.InternalError("")
	}

	if identity.Deleted {
		return utils.NotFoundError("Listing")
	}

	if identity.UserID == userID {
		return utils.AuthorizationError("owners cannot favorite their own listing")
	}

	if !identity.ActiveVersionID.Valid || identity.ActiveVersionID.Int64 == 0 {
		return utils.ValidationError("listing", "listing has no active version to favorite")
	}

	listing, listErr := ls.listingRepository.GetListingVersionByID(ctx, tx, identity.ActiveVersionID.Int64)
	if listErr != nil {
		if errors.Is(listErr, sql.ErrNoRows) {
			return utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, listErr)
		logger.Error("listing.favorite.add.listing_error", "err", listErr, "active_version_id", identity.ActiveVersionID.Int64)
		return utils.InternalError("")
	}

	if listing.Status() != listingmodel.StatusPublished {
		return utils.ValidationError("listing", "only published listings can be favorited")
	}

	if favErr := ls.favoriteRepo.Add(ctx, tx, userID, listingIdentityID); favErr != nil {
		utils.SetSpanError(ctx, favErr)
		logger.Error("listing.favorite.add.exec_error", "err", favErr, "listing_identity_id", listingIdentityID, "user_id", userID)
		return utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.favorite.add.tx_commit_error", "err", cmErr)
		return utils.InternalError("")
	}

	return nil
}

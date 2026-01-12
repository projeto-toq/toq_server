package listingservices

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// RemoveFavoriteListing unlinks the authenticated user from a listing favorite relationship.
// Operation is idempotent.
func (ls *listingService) RemoveFavoriteListing(ctx context.Context, listingIdentityID int64) (err error) {
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
		logger.Error("listing.favorite.remove.tx_start_error", "err", txErr)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.favorite.remove.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if rmErr := ls.favoriteRepo.Remove(ctx, tx, userID, listingIdentityID); rmErr != nil {
		utils.SetSpanError(ctx, rmErr)
		logger.Error("listing.favorite.remove.exec_error", "err", rmErr, "listing_identity_id", listingIdentityID, "user_id", userID)
		return utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.favorite.remove.tx_commit_error", "err", cmErr)
		return utils.InternalError("")
	}

	return nil
}

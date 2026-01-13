package listingservices

import (
	"context"
	"net/http"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// registerListingView increments the view counter for a listing identity using a short-lived write transaction.
// It returns the updated total views after the increment.
func (ls *listingService) registerListingView(ctx context.Context, listingIdentityID int64) (int64, error) {
	tx, err := ls.gsi.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("listing.detail.view.tx_start_error", "listing_identity_id", listingIdentityID, "err", err)
		return 0, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Failed to start transaction to register listing view", map[string]any{
			"stage":             "view_tx_start",
			"listingIdentityId": listingIdentityID,
		})
	}

	defer func() {
		_ = ls.gsi.RollbackTransaction(ctx, tx)
	}()

	total, incErr := ls.viewRepo.IncrementAndGet(ctx, tx, listingIdentityID)
	if incErr != nil {
		utils.SetSpanError(ctx, incErr)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("listing.detail.view.increment_error", "listing_identity_id", listingIdentityID, "err", incErr)
		return 0, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Failed to register listing view", map[string]any{
			"stage":             "view_increment",
			"listingIdentityId": listingIdentityID,
		})
	}

	if commitErr := ls.gsi.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("listing.detail.view.tx_commit_error", "listing_identity_id", listingIdentityID, "err", commitErr)
		return 0, utils.NewHTTPErrorWithSource(http.StatusInternalServerError, "Failed to commit listing view transaction", map[string]any{
			"stage":             "view_tx_commit",
			"listingIdentityId": listingIdentityID,
		})
	}

	return total, nil
}

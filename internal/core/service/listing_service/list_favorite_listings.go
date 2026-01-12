package listingservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListFavoriteListingsOutput encapsulates paginated favorite listings for the authenticated user.
type ListFavoriteListingsOutput struct {
	Items []ListListingsItem
	Total int64
	Page  int
	Limit int
}

// ListFavoriteListings returns paginated favorites for the authenticated user.
func (ls *listingService) ListFavoriteListings(ctx context.Context, page, limit int) (output ListFavoriteListingsOutput, err error) {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		return output, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	userID, err := ls.gsi.GetUserIDFromContext(ctx)
	if err != nil {
		return output, err
	}

	tx, txErr := ls.gsi.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.favorite.list.tx_start_error", "err", txErr)
		return output, utils.InternalError("")
	}
	defer func() {
		_ = ls.gsi.RollbackTransaction(ctx, tx)
	}()

	ids, total, listErr := ls.favoriteRepo.ListByUser(ctx, tx, userID, page, limit)
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		logger.Error("listing.favorite.list.repo_error", "err", listErr)
		return output, utils.InternalError("")
	}

	output.Page = page
	output.Limit = limit
	output.Total = total

	if len(ids) == 0 || total == 0 {
		output.Items = []ListListingsItem{}
		return output, nil
	}

	favoritesCount, countErr := ls.favoriteRepo.CountByListingIdentities(ctx, tx, ids)
	if countErr != nil {
		utils.SetSpanError(ctx, countErr)
		logger.Error("listing.favorite.list.count_error", "err", countErr, "ids", ids)
		return output, utils.InternalError("")
	}

	items := make([]ListListingsItem, 0, len(ids))
	for _, identityID := range ids {
		identity, idErr := ls.listingRepository.GetListingIdentityByID(ctx, tx, identityID)
		if idErr != nil {
			if !errors.Is(idErr, sql.ErrNoRows) {
				utils.SetSpanError(ctx, idErr)
				logger.Warn("listing.favorite.list.identity_error", "err", idErr, "listing_identity_id", identityID)
			}
			continue
		}

		if !identity.ActiveVersionID.Valid || identity.ActiveVersionID.Int64 == 0 {
			logger.Warn("listing.favorite.list.no_active_version", "listing_identity_id", identityID)
			continue
		}

		listing, vErr := ls.listingRepository.GetListingVersionByID(ctx, tx, identity.ActiveVersionID.Int64)
		if vErr != nil {
			if !errors.Is(vErr, sql.ErrNoRows) {
				utils.SetSpanError(ctx, vErr)
				logger.Warn("listing.favorite.list.version_error", "err", vErr, "active_version_id", identity.ActiveVersionID.Int64)
			}
			continue
		}

		// Propagate identity metadata to listing instance for response assembly
		listing.SetIdentityID(identityID)
		listing.SetUUID(identity.UUID)
		listing.SetActiveVersionID(identity.ActiveVersionID.Int64)

		// Attach draft version metadata when available (best-effort)
		if draft, dErr := ls.listingRepository.GetDraftVersionByListingIdentityID(ctx, tx, identityID); dErr == nil && draft != nil {
			listing.SetDraftVersion(draft)
		}

		items = append(items, ListListingsItem{
			Listing:        listing,
			FavoritesCount: favoritesCount[identityID],
			IsFavorite:     true,
		})
	}

	output.Items = items
	return output, nil
}

package listingservices

import (
	"context"
	"strings"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) ListCatalogValues(ctx context.Context, category string, includeInactive bool) ([]listingmodel.CatalogValueInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	category = strings.TrimSpace(strings.ToLower(category))
	if !listingmodel.IsValidCatalogCategory(category) {
		return nil, utils.ValidationError("category", "invalid catalog category")
	}

	tx, txErr := ls.gsi.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.catalog.list.tx_start_error", "err", txErr, "category", category)
		return nil, utils.InternalError("")
	}

	values, repoErr := ls.gsi.ListCatalogValues(ctx, tx, category, includeInactive)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("listing.catalog.list.query_error", "err", repoErr, "category", category, "include_inactive", includeInactive)
		_ = ls.gsi.RollbackTransaction(ctx, tx)
		return nil, utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.catalog.list.tx_commit_error", "err", cmErr, "category", category)
		_ = ls.gsi.RollbackTransaction(ctx, tx)
		return nil, utils.InternalError("")
	}

	return values, nil
}

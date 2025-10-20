package listingservices

import (
	"context"
	"database/sql"
	"strings"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetCatalogValueDetail returns a catalog value by category and identifier.
func (ls *listingService) GetCatalogValueDetail(ctx context.Context, category string, id uint8) (listingmodel.CatalogValueInterface, error) {
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

	if id == 0 {
		return nil, utils.ValidationError("id", "id must be greater than zero")
	}

	tx, txErr := ls.gsi.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.catalog.detail.tx_start_error", "err", txErr, "category", category, "id", id)
		return nil, utils.InternalError("")
	}

	value, repoErr := ls.listingRepository.GetCatalogValueByID(ctx, tx, category, id)
	if repoErr != nil {
		if repoErr == sql.ErrNoRows {
			_ = ls.gsi.RollbackTransaction(ctx, tx)
			return nil, utils.NotFoundError("catalog_value")
		}
		utils.SetSpanError(ctx, repoErr)
		logger.Error("listing.catalog.detail.repo_error", "err", repoErr, "category", category, "id", id)
		_ = ls.gsi.RollbackTransaction(ctx, tx)
		return nil, utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.catalog.detail.tx_commit_error", "err", cmErr, "category", category, "id", id)
		_ = ls.gsi.RollbackTransaction(ctx, tx)
		return nil, utils.InternalError("")
	}

	return value, nil
}

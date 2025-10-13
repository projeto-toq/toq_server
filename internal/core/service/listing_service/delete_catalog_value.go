package listingservices

import (
	"context"
	"database/sql"
	"strings"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) DeleteCatalogValue(ctx context.Context, category string, id uint8) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	category = strings.TrimSpace(strings.ToLower(category))
	if !listingmodel.IsValidCatalogCategory(category) {
		return utils.ValidationError("category", "invalid catalog category")
	}

	if id == 0 {
		return utils.ValidationError("id", "invalid catalog identifier")
	}

	tx, txErr := ls.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.catalog.delete.tx_start_error", "err", txErr, "category", category, "catalog_id", id)
		return utils.InternalError("")
	}
	rollback := true
	defer func() {
		if !rollback {
			return
		}
		if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("listing.catalog.delete.tx_rollback_error", "err", rbErr, "category", category, "catalog_id", id)
		}
	}()

	if delErr := ls.gsi.SoftDeleteCatalogValue(ctx, tx, category, id); delErr != nil {
		if delErr == sql.ErrNoRows {
			return utils.NotFoundError("catalog value")
		}
		utils.SetSpanError(ctx, delErr)
		logger.Error("listing.catalog.delete.exec_error", "err", delErr, "category", category, "catalog_id", id)
		return utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.catalog.delete.tx_commit_error", "err", cmErr, "category", category, "catalog_id", id)
		return utils.InternalError("")
	}
	rollback = false

	logger.Info("listing.catalog.deleted", "category", category, "catalog_id", id)
	return nil
}

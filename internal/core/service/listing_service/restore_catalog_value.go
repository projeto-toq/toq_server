package listingservices

import (
	"context"
	"database/sql"
	"strings"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) RestoreCatalogValue(ctx context.Context, input RestoreCatalogValueInput) (listingmodel.CatalogValueInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	category := strings.TrimSpace(strings.ToLower(input.Category))
	if !listingmodel.IsValidCatalogCategory(category) {
		return nil, utils.ValidationError("category", "invalid catalog category")
	}

	if input.ID == 0 {
		return nil, utils.ValidationError("id", "invalid catalog identifier")
	}

	tx, txErr := ls.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.catalog.restore.tx_start_error", "err", txErr, "category", category, "catalog_id", input.ID)
		return nil, utils.InternalError("")
	}
	rollback := true
	defer func() {
		if !rollback {
			return
		}
		if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("listing.catalog.restore.tx_rollback_error", "err", rbErr, "category", category, "catalog_id", input.ID)
		}
	}()

	value, getErr := ls.listingRepository.GetCatalogValueByID(ctx, tx, category, input.ID)
	if getErr != nil {
		if getErr == sql.ErrNoRows {
			return nil, utils.NotFoundError("catalog value")
		}
		utils.SetSpanError(ctx, getErr)
		logger.Error("listing.catalog.restore.get_error", "err", getErr, "category", category, "catalog_id", input.ID)
		return nil, utils.InternalError("")
	}

	if value.IsActive() {
		return nil, utils.ConflictError("catalog value already active")
	}

	value.SetIsActive(true)

	if updateErr := ls.listingRepository.UpdateCatalogValue(ctx, tx, value); updateErr != nil {
		if updateErr == sql.ErrNoRows {
			return nil, utils.NotFoundError("catalog value")
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("listing.catalog.restore.exec_error", "err", updateErr, "category", category, "catalog_id", input.ID)
		return nil, utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.catalog.restore.tx_commit_error", "err", cmErr, "category", category, "catalog_id", input.ID)
		return nil, utils.InternalError("")
	}
	rollback = false

	logger.Info("listing.catalog.restored", "category", category, "catalog_id", value.ID())
	return value, nil
}

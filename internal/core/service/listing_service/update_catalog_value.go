package listingservices

import (
	"context"
	"database/sql"
	"strings"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) UpdateCatalogValue(ctx context.Context, input UpdateCatalogValueInput) (listingmodel.CatalogValueInterface, error) {
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
		logger.Error("listing.catalog.update.tx_start_error", "err", txErr, "category", category, "catalog_id", input.ID)
		return nil, utils.InternalError("")
	}
	rollback := true
	defer func() {
		if !rollback {
			return
		}
		if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("listing.catalog.update.tx_rollback_error", "err", rbErr, "category", category, "catalog_id", input.ID)
		}
	}()

	value, getErr := ls.listingRepository.GetCatalogValueByID(ctx, tx, category, input.ID)
	if getErr != nil {
		if getErr == sql.ErrNoRows {
			return nil, utils.NotFoundError("catalog value")
		}
		utils.SetSpanError(ctx, getErr)
		logger.Error("listing.catalog.update.get_error", "err", getErr, "category", category, "catalog_id", input.ID)
		return nil, utils.InternalError("")
	}

	if input.Slug.IsPresent() {
		if input.Slug.IsNull() {
			return nil, utils.ValidationError("slug", "slug cannot be null")
		}
		if newSlug, ok := input.Slug.Value(); ok {
			normalized := strings.TrimSpace(strings.ToLower(newSlug))
			if normalized == "" {
				return nil, utils.ValidationError("slug", "slug is required")
			}
			if normalized != value.Slug() {
				if existing, slugErr := ls.listingRepository.GetCatalogValueBySlug(ctx, tx, category, normalized); slugErr == nil {
					if existing.ID() != value.ID() {
						return nil, utils.ConflictError("catalog slug already exists")
					}
				} else if slugErr != sql.ErrNoRows {
					utils.SetSpanError(ctx, slugErr)
					logger.Error("listing.catalog.update.slug_lookup_error", "err", slugErr, "category", category, "catalog_id", input.ID, "slug", normalized)
					return nil, utils.InternalError("")
				}
				value.SetSlug(normalized)
			}
		}
	}

	if input.Label.IsPresent() {
		if input.Label.IsNull() {
			return nil, utils.ValidationError("label", "label cannot be null")
		}
		if newLabel, ok := input.Label.Value(); ok {
			normalized := strings.TrimSpace(newLabel)
			if normalized == "" {
				return nil, utils.ValidationError("label", "label is required")
			}
			value.SetLabel(normalized)
		}
	}

	if input.Description.IsPresent() {
		if input.Description.IsNull() {
			value.SetDescription(nil)
		} else if desc, ok := input.Description.Value(); ok {
			copyDesc := strings.TrimSpace(desc)
			if copyDesc == "" {
				value.SetDescription(nil)
			} else {
				value.SetDescription(&copyDesc)
			}
		}
	}

	if input.IsActive.IsPresent() {
		if input.IsActive.IsNull() {
			value.SetIsActive(false)
		} else if active, ok := input.IsActive.Value(); ok {
			value.SetIsActive(active)
		}
	}

	if updateErr := ls.listingRepository.UpdateCatalogValue(ctx, tx, value); updateErr != nil {
		if updateErr == sql.ErrNoRows {
			return nil, utils.NotFoundError("catalog value")
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("listing.catalog.update.exec_error", "err", updateErr, "category", category, "catalog_id", input.ID)
		return nil, utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.catalog.update.tx_commit_error", "err", cmErr, "category", category, "catalog_id", input.ID)
		return nil, utils.InternalError("")
	}
	rollback = false

	logger.Info("listing.catalog.updated", "category", category, "catalog_id", value.ID())
	return value, nil
}

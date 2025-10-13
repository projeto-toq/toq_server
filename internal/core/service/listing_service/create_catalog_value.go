package listingservices

import (
	"context"
	"database/sql"
	"strings"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) CreateCatalogValue(ctx context.Context, input CreateCatalogValueInput) (listingmodel.CatalogValueInterface, error) {
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

	slug := strings.TrimSpace(strings.ToLower(input.Slug))
	if slug == "" {
		return nil, utils.ValidationError("slug", "slug is required")
	}

	label := strings.TrimSpace(input.Label)
	if label == "" {
		return nil, utils.ValidationError("label", "label is required")
	}

	input.Category = category
	input.Slug = slug
	input.Label = label

	tx, txErr := ls.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.catalog.create.tx_start_error", "err", txErr, "category", category, "slug", slug)
		return nil, utils.InternalError("")
	}
	rollback := true
	defer func() {
		if !rollback {
			return
		}
		if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("listing.catalog.create.tx_rollback_error", "err", rbErr, "category", category, "slug", slug)
		}
	}()

	if _, dupErr := ls.listingRepository.GetCatalogValueBySlug(ctx, tx, category, slug); dupErr == nil {
		return nil, utils.ConflictError("catalog slug already exists")
	} else if dupErr != sql.ErrNoRows {
		utils.SetSpanError(ctx, dupErr)
		logger.Error("listing.catalog.create.slug_lookup_error", "err", dupErr, "category", category, "slug", slug)
		return nil, utils.InternalError("")
	}

	nextID, idErr := ls.listingRepository.GetNextCatalogValueID(ctx, tx, category)
	if idErr != nil {
		utils.SetSpanError(ctx, idErr)
		logger.Error("listing.catalog.create.next_id_error", "err", idErr, "category", category)
		return nil, utils.InternalError("")
	}

	value := input.ToDomain(nextID)

	if createErr := ls.listingRepository.CreateCatalogValue(ctx, tx, value); createErr != nil {
		utils.SetSpanError(ctx, createErr)
		logger.Error("listing.catalog.create.exec_error", "err", createErr, "category", category, "slug", slug)
		return nil, utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.catalog.create.tx_commit_error", "err", cmErr, "category", category, "slug", slug)
		return nil, utils.InternalError("")
	}
	rollback = false

	logger.Info("listing.catalog.created", "category", category, "catalog_id", value.ID(), "slug", value.Slug())
	return value, nil
}

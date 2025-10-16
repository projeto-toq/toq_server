package listingservices

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) resolveCatalogValue(ctx context.Context, tx *sql.Tx, category, field string, selection CatalogSelection) (uint8, error) {
	if selection.HasID() {
		value, err := ls.listingRepository.GetCatalogValueByID(ctx, tx, category, selection.IDValue())
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, utils.ValidationError(field, "invalid value")
			}
			utils.SetSpanError(ctx, err)
			return 0, utils.InternalError("")
		}
		if !value.IsActive() {
			return 0, utils.ValidationError(field, "value is inactive")
		}
		return value.NumericValue(), nil
	}

	slug := selection.SlugValue()
	if slug == "" {
		return 0, utils.ValidationError(field, "invalid value")
	}

	value, err := ls.listingRepository.GetCatalogValueBySlug(ctx, tx, category, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, utils.ValidationError(field, "invalid value")
		}
		utils.SetSpanError(ctx, err)
		return 0, utils.InternalError("")
	}

	if !value.IsActive() {
		return 0, utils.ValidationError(field, "value is inactive")
	}

	return value.NumericValue(), nil
}

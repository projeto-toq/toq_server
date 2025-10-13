package listingservices

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) ensureCatalogValue(ctx context.Context, tx *sql.Tx, category string, id uint8, field string) error {
	if id == 0 {
		return nil
	}

	value, err := ls.gsi.GetCatalogValueByID(ctx, tx, category, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ValidationError(field, "invalid value")
		}
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	if !value.IsActive() {
		return utils.ValidationError(field, "value is inactive")
	}

	return nil
}

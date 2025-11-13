package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) CreateFeature(ctx context.Context, tx *sql.Tx, feature listingmodel.FeatureInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	statement := `INSERT INTO features (listing_version_id, feature_id, qty) VALUES (?, ?, ?);`

	result, execErr := la.ExecContext(ctx, tx, "insert", statement, feature.ListingVersionID(), feature.FeatureID(), feature.Quantity())
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.create_feature.exec_error", "error", execErr)
		return fmt.Errorf("exec create feature: %w", execErr)
	}

	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.listing.create_feature.last_insert_error", "error", lastErr)
		return fmt.Errorf("last insert id for feature: %w", lastErr)
	}

	feature.SetID(id)

	return nil
}

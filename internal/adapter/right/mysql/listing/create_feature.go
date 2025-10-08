package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) CreateFeature(ctx context.Context, tx *sql.Tx, feature listingmodel.FeatureInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	sql := `INSERT INTO features (listing_id, feature_id, qty) VALUES (?, ?, ?);`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_feature.prepare_error", "error", err)
		return fmt.Errorf("prepare create feature: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, feature.ListingID(), feature.FeatureID(), feature.Quantity())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_feature.exec_error", "error", err)
		return fmt.Errorf("exec create feature: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_feature.last_insert_error", "error", err)
		return fmt.Errorf("last insert id for feature: %w", err)
	}

	feature.SetID(id)

	return nil
}

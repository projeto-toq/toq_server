package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/entity"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetBaseFeatures(ctx context.Context, tx *sql.Tx) (features []listingmodel.BaseFeatureInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	sql := `SELECT * FROM base_features ORDER BY priority ASC;`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_base_features.prepare_error", "error", err)
		return nil, fmt.Errorf("prepare get base features: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_base_features.query_error", "error", err)
		return nil, fmt.Errorf("query get base features: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		entity := listingentity.EntityBaseFeature{}
		err = rows.Scan(
			&entity.ID,
			&entity.Feature,
			&entity.Description,
			&entity.Priority,
		)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.get_base_features.scan_error", "error", err)
			return nil, fmt.Errorf("scan base feature row: %w", err)
		}
		feature := listingmodel.NewBaseFeature()
		feature.SetID(entity.ID)
		feature.SetFeature(entity.Feature)
		feature.SetDescription(entity.Description)
		feature.SetPriority(entity.Priority)

		features = append(features, feature)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_base_features.rows_error", "error", err)
		return nil, fmt.Errorf("rows iteration for base features: %w", err)
	}

	return features, nil
}

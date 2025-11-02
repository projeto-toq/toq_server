package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/entity"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetEntityFeaturesByListing(ctx context.Context, tx *sql.Tx, listingID int64) (features []listingentity.EntityFeature, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT * FROM features WHERE listing_id = ?;`

	rows, queryErr := la.QueryContext(ctx, tx, "select", query, listingID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.listing.get_entity_features.query_error", "error", queryErr)
		return nil, fmt.Errorf("query features by listing: %w", queryErr)
	}
	defer rows.Close()

	for rows.Next() {
		feature := listingentity.EntityFeature{}
		err = rows.Scan(
			&feature.ID,
			&feature.ListingID,
			&feature.FeatureID,
			&feature.Quantity,
		)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.get_entity_features.scan_error", "error", err)
			return nil, fmt.Errorf("scan feature row: %w", err)
		}

		features = append(features, feature)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_features.rows_error", "error", err)
		return nil, fmt.Errorf("rows iteration for features: %w", err)
	}

	return features, nil
}

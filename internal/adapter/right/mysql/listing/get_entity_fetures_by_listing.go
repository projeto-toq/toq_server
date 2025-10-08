package mysqllistingadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	listingentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing/entity"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_features.prepare_error", "error", err)
		return nil, fmt.Errorf("prepare get features: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, listingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_features.query_error", "error", err)
		return nil, fmt.Errorf("query features by listing: %w", err)
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

package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	listingentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing/entity"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetEntityFeaturesByListing(ctx context.Context, tx *sql.Tx, listingID int64) (features []listingentity.EntityFeature, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM features WHERE listing_id = ?;`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("Error preparing statement on mysqllistingadapter/GetEntityFeaturesByListing", "error", err)
		err = fmt.Errorf("prepare get features: %w", err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, listingID)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("Error executing query on mysqllistingadapter/GetEntityFeaturesByListing", "error", err)
		err = fmt.Errorf("query features by listing: %w", err)
		return
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
			slog.Error("Error scanning row on mysqllistingadapter/GetEntityFeaturesByListing", "error", err)
			err = fmt.Errorf("scan feature row: %w", err)
			return
		}

		features = append(features, feature)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating over rows on mysqllistingadapter/GetEntityFeaturesByListing", "error", err)
		err = fmt.Errorf("rows iteration for features: %w", err)
		return
	}

	return
}

package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	listingentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing/entity"
	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetBaseFeatures(ctx context.Context, tx *sql.Tx) (features []listingmodel.BaseFeatureInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `SELECT * FROM base_features ORDER BY priority ASC;`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		slog.Error("Error preparing statement on mysqllistingadapter/GetBasefeatures", "error", err)
		err = fmt.Errorf("prepare get base features: %w", err)
		return
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		slog.Error("Error executing statement on mysqllistingadapter/GetBasefeatures", "error", err)
		err = fmt.Errorf("query get base features: %w", err)
		return
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
			slog.Error("Error scanning row on mysqllistingadapter/GetBasefeatures", "error", err)
			err = fmt.Errorf("scan base feature row: %w", err)
			return
		}
		feature := listingmodel.NewBaseFeature()
		feature.SetID(entity.ID)
		feature.SetFeature(entity.Feature)
		feature.SetDescription(entity.Description)
		feature.SetPriority(entity.Priority)

		features = append(features, feature)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating over rows on mysqllistingadapter/GetBasefeatures", "error", err)
		err = fmt.Errorf("rows iteration for base features: %w", err)
		return
	}

	return
}

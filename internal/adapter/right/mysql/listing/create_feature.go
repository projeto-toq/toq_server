package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) CreateFeature(ctx context.Context, tx *sql.Tx, feature listingmodel.FeatureInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `INSERT INTO features (listing_id, feature_id, qty) VALUES (?, ?, ?);`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		slog.Error("mysqllistingadapter/CreateFeature: error preparing statement", "error", err)
		err = fmt.Errorf("prepare create feature: %w", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, feature.ListingID(), feature.FeatureID(), feature.Quantity())
	if err != nil {
		slog.Error("mysqllistingadapter/CreateFeature: error executing statement", "error", err)
		err = fmt.Errorf("exec create feature: %w", err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("mysqllistingadapter/CreateFeature: error getting last insert ID", "error", err)
		err = fmt.Errorf("last insert id for feature: %w", err)
		return
	}

	feature.SetID(id)

	return
}

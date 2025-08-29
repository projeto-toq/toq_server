package mysqllistingadapter

import (
	"context"
	"database/sql"
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
		err = utils.ErrInternalServer
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, feature.ListingID(), feature.FeatureID(), feature.Quantity())
	if err != nil {
		slog.Error("mysqllistingadapter/CreateFeature: error executing statement", "error", err)
		err = utils.ErrInternalServer
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("mysqllistingadapter/CreateFeature: error getting last insert ID", "error", err)
		err = utils.ErrInternalServer
		return
	}

	feature.SetID(id)

	return
}

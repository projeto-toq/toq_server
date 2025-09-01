package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) CreateExchangePlace(ctx context.Context, tx *sql.Tx, place listingmodel.ExchangePlaceInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `INSERT INTO exchange_places (listing_id, neighborhood, city, state) VALUES (?, ?, ?, ?);`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		slog.Error("mysqllistingadapter/CreateExchangePlace: error preparing statement", "error", err)
		err = fmt.Errorf("prepare create exchange place: %w", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, place.ListingID(), place.Neighborhood(), place.City(), place.State())
	if err != nil {
		slog.Error("mysqllistingadapter/CreateExchangePlace: error executing statement", "error", err)
		err = fmt.Errorf("exec create exchange place: %w", err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("mysqllistingadapter/CreateExchangePlace: error getting last insert ID", "error", err)
		err = fmt.Errorf("last insert id for exchange place: %w", err)
		return
	}

	place.SetID(id)

	return
}

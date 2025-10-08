package mysqllistingadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	listingentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing/entity"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetEntityExchangePlacesByListing(ctx context.Context, tx *sql.Tx, listingID int64) (places []listingentity.EntityExchangePlace, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT * FROM exchange_places WHERE listing_id = ?;`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_exchange_places.prepare_error", "error", err)
		return nil, fmt.Errorf("prepare get exchange places: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, listingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_exchange_places.query_error", "error", err)
		return nil, fmt.Errorf("query exchange places by listing: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		place := listingentity.EntityExchangePlace{}
		err = rows.Scan(
			&place.ID,
			&place.ListingID,
			&place.Neighborhood,
			&place.City,
			&place.State,
		)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.get_entity_exchange_places.scan_error", "error", err)
			return nil, fmt.Errorf("scan exchange place row: %w", err)
		}

		places = append(places, place)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_exchange_places.rows_error", "error", err)
		return nil, fmt.Errorf("rows iteration for exchange places: %w", err)
	}

	return places, nil
}

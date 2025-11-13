package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/entity"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetEntityExchangePlacesByListing(ctx context.Context, tx *sql.Tx, listingVersionID int64) (places []listingentity.EntityExchangePlace, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT * FROM exchange_places WHERE listing_version_id = ?;`

	rows, queryErr := la.QueryContext(ctx, tx, "select", query, listingVersionID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.listing.get_entity_exchange_places.query_error", "error", queryErr)
		return nil, fmt.Errorf("query exchange places by listing: %w", queryErr)
	}
	defer rows.Close()

	for rows.Next() {
		place := listingentity.EntityExchangePlace{}
		err = rows.Scan(
			&place.ID,
			&place.ListingVersionID,
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

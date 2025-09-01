package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	listingentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing/entity"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetEntityExchangePlacesByListing(ctx context.Context, tx *sql.Tx, listingID int64) (places []listingentity.EntityExchangePlace, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM exchange_places WHERE listing_id = ?;`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("Error preparing statement in GetEntityExchangePlacesByListing", "error", err)
		err = fmt.Errorf("prepare get exchange places: %w", err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, listingID)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("Error executing query in GetEntityExchangePlacesByListing", "error", err)
		err = fmt.Errorf("query exchange places by listing: %w", err)
		return
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
			slog.Error("Error scanning row in GetEntityExchangePlacesByListing", "error", err)
			err = fmt.Errorf("scan exchange place row: %w", err)
			return
		}

		places = append(places, place)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating over rows in GetEntityExchangePlacesByListing", "error", err)
		err = fmt.Errorf("rows iteration for exchange places: %w", err)
		return
	}

	return
}

package mysqllistingadapter

import (
	"context"
	"database/sql"
	"log/slog"

	listingentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing/entity"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetEntityFinancingBlockersByListing(ctx context.Context, tx *sql.Tx, listingID int64) (blockers []listingentity.EntityFinancingBlocker, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM financing_blockers WHERE listing_id = ?;`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("Error preparing statement on mysqllistingadapter/GetEntityFinancingBlockerByListing", "error", err)
		err = utils.ErrInternalServer
		return
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, listingID)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("Error executing query on mysqllistingadapter/GetEntityFinancingBlockerByListing", "error", err)
		err = utils.ErrInternalServer
		return
	}
	defer rows.Close()

	for rows.Next() {
		blocker := listingentity.EntityFinancingBlocker{}
		err = rows.Scan(
			&blocker.ID,
			&blocker.ListingID,
			&blocker.Blocker,
		)
		if err != nil {
			slog.Error("Error scanning row on mysqllistingadapter/GetEntityFinancingBlockerByListing", "error", err)
			err = utils.ErrInternalServer
			return
		}

		blockers = append(blockers, blocker)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating over rows on mysqllistingadapter/GetEntityFinancingBlockerByListing", "error", err)
		err = utils.ErrInternalServer
		return
	}

	return
}

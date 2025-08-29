package mysqllistingadapter

import (
	"context"
	"database/sql"
	"log/slog"

	listingentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing/entity"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetEntityGuaranteesByListing(ctx context.Context, tx *sql.Tx, listingID int64) (guarantees []listingentity.EntityGuarantee, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM guarantees WHERE listing_id = ?;`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("Error preparing statement in GetEntityGuaranteesByListing", "error", err)
		err = utils.ErrInternalServer
		return
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, listingID)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("Error executing query in GetEntityGuaranteesByListing", "error", err)
		err = utils.ErrInternalServer
		return
	}
	defer rows.Close()

	for rows.Next() {
		guarantee := listingentity.EntityGuarantee{}
		err = rows.Scan(
			&guarantee.ID,
			&guarantee.ListingID,
			&guarantee.Priority,
			&guarantee.Guarantee,
		)
		if err != nil {
			slog.Error("Error scanning row in GetEntityGuaranteesByListing", "error", err)
			err = utils.ErrInternalServer
			return
		}

		guarantees = append(guarantees, guarantee)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating over rows in GetEntityGuaranteesByListing", "error", err)
		err = utils.ErrInternalServer
		return
	}

	return
}

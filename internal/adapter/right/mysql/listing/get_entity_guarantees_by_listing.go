package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/entity"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetEntityGuaranteesByListing(ctx context.Context, tx *sql.Tx, listingID int64) (guarantees []listingentity.EntityGuarantee, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT * FROM guarantees WHERE listing_id = ?;`

	rows, queryErr := la.QueryContext(ctx, tx, "select", query, listingID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.listing.get_entity_guarantees.query_error", "error", queryErr)
		return nil, fmt.Errorf("query guarantees by listing: %w", queryErr)
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
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.get_entity_guarantees.scan_error", "error", err)
			return nil, fmt.Errorf("scan guarantee row: %w", err)
		}

		guarantees = append(guarantees, guarantee)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_guarantees.rows_error", "error", err)
		return nil, fmt.Errorf("rows iteration for guarantees: %w", err)
	}

	return guarantees, nil
}

package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/listing/entity"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetEntityFinancingBlockersByListing(ctx context.Context, tx *sql.Tx, listingID int64) (blockers []listingentity.EntityFinancingBlocker, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT * FROM financing_blockers WHERE listing_id = ?;`

	rows, queryErr := la.QueryContext(ctx, tx, "select", query, listingID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.listing.get_entity_financing_blockers.query_error", "error", queryErr)
		return nil, fmt.Errorf("query financing blockers by listing: %w", queryErr)
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
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.listing.get_entity_financing_blockers.scan_error", "error", err)
			return nil, fmt.Errorf("scan financing blocker row: %w", err)
		}

		blockers = append(blockers, blocker)
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_financing_blockers.rows_error", "error", err)
		return nil, fmt.Errorf("rows iteration for financing blockers: %w", err)
	}

	return blockers, nil
}

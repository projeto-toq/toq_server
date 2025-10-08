package mysqllistingadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	listingentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/listing/entity"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_financing_blockers.prepare_error", "error", err)
		return nil, fmt.Errorf("prepare get financing blockers: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, listingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_financing_blockers.query_error", "error", err)
		return nil, fmt.Errorf("query financing blockers by listing: %w", err)
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

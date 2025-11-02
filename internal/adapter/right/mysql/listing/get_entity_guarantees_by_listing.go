package mysqllistingadapter

import (
	"context"
	"database/sql"
	"errors"
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
	defer la.ObserveOnComplete("select", query)()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_guarantees.prepare_error", "error", err)
		return nil, fmt.Errorf("prepare get guarantees: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, listingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_entity_guarantees.query_error", "error", err)
		return nil, fmt.Errorf("query guarantees by listing: %w", err)
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

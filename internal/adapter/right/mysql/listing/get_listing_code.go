package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetListingCode(ctx context.Context, tx *sql.Tx) (code uint32, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE listing_sequence SET id=LAST_INSERT_ID(id+1);`
	defer la.ObserveOnComplete("update", query)()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_listing_code.prepare_error", "error", err)
		return 0, fmt.Errorf("prepare get listing code: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_listing_code.exec_error", "error", err)
		return 0, fmt.Errorf("exec get listing code: %w", err)
	}

	code64, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.get_listing_code.last_insert_error", "error", err)
		return 0, fmt.Errorf("last insert id for listing code: %w", err)
	}

	code = uint32(code64)

	return code, nil
}

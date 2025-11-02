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

	result, execErr := la.ExecContext(ctx, tx, "update", query)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.get_listing_code.exec_error", "error", execErr)
		return 0, fmt.Errorf("exec get listing code: %w", execErr)
	}

	code64, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.listing.get_listing_code.last_insert_error", "error", lastErr)
		return 0, fmt.Errorf("last insert id for listing code: %w", lastErr)
	}

	code = uint32(code64)

	return code, nil
}

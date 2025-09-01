package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetListingCode(ctx context.Context, tx *sql.Tx) (code uint32, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `UPDATE listing_sequence SET id=LAST_INSERT_ID(id+1);`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		slog.Error("Error preparing statement on msqllistingadapter/GetListingCode", "error", err)
		err = fmt.Errorf("prepare get listing code: %w", err)
		return
	}

	result, err := stmt.ExecContext(ctx)
	if err != nil {
		slog.Error("Error executing statement on msqllistingadapter/GetListingCode", "error", err)
		err = fmt.Errorf("exec get listing code: %w", err)
		return
	}

	code64, err := result.LastInsertId()
	if err != nil {
		slog.Error("Error getting last insert id on msqllistingadapter/GetListingCode", "error", err)
		err = fmt.Errorf("last insert id for listing code: %w", err)
		return
	}

	code = uint32(code64)

	return
}

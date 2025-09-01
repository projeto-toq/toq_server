package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) CreateFinancingBlocker(ctx context.Context, tx *sql.Tx, blocker listingmodel.FinancingBlockerInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `INSERT INTO financing_blockers (listing_id, blocker) VALUES (?, ?);`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		slog.Error("mysqllistingadapter/CreateFinancingBlocker: error preparing statement", "error", err)
		err = fmt.Errorf("prepare create financing blocker: %w", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, blocker.ListingID(), blocker.Blocker())
	if err != nil {
		slog.Error("mysqllistingadapter/CreateFinancingBlocker: error executing statement", "error", err)
		err = fmt.Errorf("exec create financing blocker: %w", err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("mysqllistingadapter/CreateFinancingBlocker: error getting last insert ID", "error", err)
		err = fmt.Errorf("last insert id for financing blocker: %w", err)
		return
	}

	blocker.SetID(id)

	return
}

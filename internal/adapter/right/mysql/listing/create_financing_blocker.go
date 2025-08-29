package mysqllistingadapter

import (
	"context"
	"database/sql"
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
		err = utils.ErrInternalServer
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, blocker.ListingID(), blocker.Blocker())
	if err != nil {
		slog.Error("mysqllistingadapter/CreateFinancingBlocker: error executing statement", "error", err)
		err = utils.ErrInternalServer
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("mysqllistingadapter/CreateFinancingBlocker: error getting last insert ID", "error", err)
		err = utils.ErrInternalServer
		return
	}

	blocker.SetID(id)

	return
}

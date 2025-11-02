package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) CreateFinancingBlocker(ctx context.Context, tx *sql.Tx, blocker listingmodel.FinancingBlockerInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	statement := `INSERT INTO financing_blockers (listing_id, blocker) VALUES (?, ?);`
	defer la.ObserveOnComplete("insert", statement)()

	stmt, err := tx.PrepareContext(ctx, statement)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_financing_blocker.prepare_error", "error", err)
		return fmt.Errorf("prepare create financing blocker: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, blocker.ListingID(), blocker.Blocker())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_financing_blocker.exec_error", "error", err)
		return fmt.Errorf("exec create financing blocker: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_financing_blocker.last_insert_error", "error", err)
		return fmt.Errorf("last insert id for financing blocker: %w", err)
	}

	blocker.SetID(id)

	return nil
}

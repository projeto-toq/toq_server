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

	statement := `INSERT INTO financing_blockers (listing_version_id, blocker) VALUES (?, ?);`

	result, execErr := la.ExecContext(ctx, tx, "insert", statement, blocker.ListingVersionID(), blocker.Blocker())
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.create_financing_blocker.exec_error", "error", execErr)
		return fmt.Errorf("exec create financing blocker: %w", execErr)
	}

	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.listing.create_financing_blocker.last_insert_error", "error", lastErr)
		return fmt.Errorf("last insert id for financing blocker: %w", lastErr)
	}

	blocker.SetID(id)

	return nil
}

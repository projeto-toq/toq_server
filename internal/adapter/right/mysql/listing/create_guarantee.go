package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) CreateGuarantee(ctx context.Context, tx *sql.Tx, guarantee listingmodel.GuaranteeInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	statement := `INSERT INTO guarantees (listing_version_id, priority, guarantee) VALUES (?, ?, ?);`

	result, execErr := la.ExecContext(ctx, tx, "insert", statement, guarantee.ListingVersionID(), guarantee.Priority(), guarantee.Guarantee())
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.create_guarantee.exec_error", "error", execErr)
		return fmt.Errorf("exec create guarantee: %w", execErr)
	}

	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.listing.create_guarantee.last_insert_error", "error", lastErr)
		return fmt.Errorf("last insert id for guarantee: %w", lastErr)
	}

	guarantee.SetID(id)

	return nil
}

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

	sql := `INSERT INTO guarantees (listing_id, priority, guarantee) VALUES (?, ?, ?);`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_guarantee.prepare_error", "error", err)
		return fmt.Errorf("prepare create guarantee: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, guarantee.ListingID(), guarantee.Priority(), guarantee.Guarantee())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_guarantee.exec_error", "error", err)
		return fmt.Errorf("exec create guarantee: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing.create_guarantee.last_insert_error", "error", err)
		return fmt.Errorf("last insert id for guarantee: %w", err)
	}

	guarantee.SetID(id)

	return nil
}

package mysqlvisitadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/converters"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *VisitAdapter) GetVisitByID(ctx context.Context, tx *sql.Tx, id int64) (listingmodel.VisitInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, listing_id, owner_id, realtor_id, scheduled_start, scheduled_end, status, cancel_reason, notes, created_by, updated_by FROM listing_visits WHERE id = ?`
	row := a.QueryRowContext(ctx, tx, "get_visit_by_id", query, id)

	visitEntity, err := scanVisitEntity(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.get.scan_error", "visit_id", id, "err", err)
		return nil, fmt.Errorf("scan visit: %w", err)
	}

	return converters.ToVisitModel(visitEntity), nil
}

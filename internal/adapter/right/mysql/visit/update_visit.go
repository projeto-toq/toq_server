package mysqlvisitadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/converters"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *VisitAdapter) UpdateVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.VisitInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToVisitEntity(visit)

	query := `UPDATE listing_visits SET listing_id = ?, owner_id = ?, realtor_id = ?, scheduled_start = ?, scheduled_end = ?, status = ?, cancel_reason = ?, notes = ?, updated_by = ? WHERE id = ?`
	result, err := a.ExecContext(ctx, tx, "update_visit", query,
		entity.ListingID,
		entity.OwnerID,
		entity.RealtorID,
		entity.ScheduledStart,
		entity.ScheduledEnd,
		entity.Status,
		entity.CancelReason,
		entity.Notes,
		entity.UpdatedBy,
		entity.ID,
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.update.exec_error", "visit_id", entity.ID, "err", err)
		return fmt.Errorf("update visit: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.update.rows_error", "visit_id", entity.ID, "err", err)
		return fmt.Errorf("visit rows affected: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

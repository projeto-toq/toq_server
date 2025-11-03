package mysqlvisitadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/converters"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *VisitAdapter) InsertVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.VisitInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToVisitEntity(visit)

	query := `INSERT INTO listing_visits (listing_id, owner_id, realtor_id, scheduled_start, scheduled_end, status, cancel_reason, notes, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := a.ExecContext(ctx, tx, "insert_visit", query,
		entity.ListingID,
		entity.OwnerID,
		entity.RealtorID,
		entity.ScheduledStart,
		entity.ScheduledEnd,
		entity.Status,
		entity.CancelReason,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.insert.exec_error", "listing_id", entity.ListingID, "err", err)
		return 0, fmt.Errorf("insert visit: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.insert.last_id_error", "listing_id", entity.ListingID, "err", err)
		return 0, fmt.Errorf("visit last insert id: %w", err)
	}

	visit.SetID(id)
	return id, nil
}

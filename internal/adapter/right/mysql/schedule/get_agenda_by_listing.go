package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	scheduleconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	scheduleentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entities"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetAgendaByListingIdentityID fetches the agenda for a given listing_identity_id from listing_agendas.
//
// Parameters:
//   - ctx: request-scoped context used for tracing and structured logging.
//   - tx: optional transaction; when provided the query must run inside it to keep consistency with sibling operations.
//   - listingIdentityID: target listing_identity_id foreign key.
//
// Returns:
//   - AgendaInterface: populated with id, listingIdentityID, ownerID and timezone.
//   - error: sql.ErrNoRows when no agenda exists; driver errors for query/scan issues.
//
// Observability:
//   - Starts a tracer span (GenerateTracer) and ensures span is ended.
//   - Uses ContextWithLogger for correlation ids; span error marked on infra failures.
func (a *ScheduleAdapter) GetAgendaByListingIdentityID(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (schedulemodel.AgendaInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, listing_identity_id, owner_id, timezone FROM listing_agendas WHERE listing_identity_id = ? LIMIT 1`

	row := a.QueryRowContext(ctx, tx, "select", query, listingIdentityID)

	var agendaEntity scheduleentity.AgendaEntity
	if err = row.Scan(&agendaEntity.ID, &agendaEntity.ListingIdentityID, &agendaEntity.OwnerID, &agendaEntity.Timezone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.get_agenda.scan_error", "listing_identity_id", listingIdentityID, "err", err)
		return nil, fmt.Errorf("scan agenda by listing id: %w", err)
	}

	return scheduleconverters.AgendaEntityToDomain(agendaEntity), nil
}

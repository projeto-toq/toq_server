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

// GetEntryByVisitID returns the agenda entry associated with a visit id.
//
// Parameters:
//   - ctx: request-scoped context for tracing/logging.
//   - tx: required transaction when consistency is needed.
//   - visitID: visit identifier to search for.
//
// Returns: AgendaEntryInterface or sql.ErrNoRows when no entry is linked to the visit; infra errors are bubbled.
// Observability: tracer span, logger propagation, span error marking on infra failures.
func (a *ScheduleAdapter) GetEntryByVisitID(ctx context.Context, tx *sql.Tx, visitID uint64) (schedulemodel.AgendaEntryInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, agenda_id, entry_type, starts_at, ends_at, blocking, reason, visit_id, photo_booking_id FROM listing_agenda_entries WHERE visit_id = ? LIMIT 1`
	row := a.QueryRowContext(ctx, tx, "select", query, visitID)

	var entryEntity scheduleentity.EntryEntity
	if err = row.Scan(&entryEntity.ID, &entryEntity.AgendaID, &entryEntity.EntryType, &entryEntity.StartsAt, &entryEntity.EndsAt, &entryEntity.Blocking, &entryEntity.Reason, &entryEntity.VisitID, &entryEntity.PhotoBookingID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.get_entry_by_visit.scan_error", "visit_id", visitID, "err", err)
		return nil, fmt.Errorf("scan agenda entry by visit: %w", err)
	}

	return scheduleconverters.EntryEntityToDomain(entryEntity), nil
}

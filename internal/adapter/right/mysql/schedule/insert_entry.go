package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	scheduleconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// InsertEntry creates a new agenda entry row and returns the generated ID.
//
// Parameters:
//   - ctx: request-scoped context for tracing/logging.
//   - tx: required transaction to keep atomicity with related writes.
//   - entry: domain entry; generated ID is set back on this object when successful.
//
// Returns: generated ID or infrastructure errors; sql.ErrNoRows is not expected for inserts.
// Observability: tracer span, logger propagation, span error marking on infra failures.
func (a *ScheduleAdapter) InsertEntry(ctx context.Context, tx *sql.Tx, entry schedulemodel.AgendaEntryInterface) (uint64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := scheduleconverters.EntryDomainToEntity(entry)

	query := `INSERT INTO listing_agenda_entries (agenda_id, entry_type, starts_at, ends_at, blocking, reason, visit_id, photo_booking_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, execErr := a.ExecContext(ctx, tx, "insert", query, entity.AgendaID, entity.EntryType, entity.StartsAt, entity.EndsAt, entity.Blocking, entity.Reason, entity.VisitID, entity.PhotoBookingID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.schedule.insert_entry.exec_error", "agenda_id", entity.AgendaID, "err", execErr)
		return 0, fmt.Errorf("insert agenda entry: %w", execErr)
	}

	id, lastIDErr := result.LastInsertId()
	if lastIDErr != nil {
		utils.SetSpanError(ctx, lastIDErr)
		logger.Error("mysql.schedule.insert_entry.last_id_error", "agenda_id", entity.AgendaID, "err", lastIDErr)
		return 0, fmt.Errorf("agenda entry last insert id: %w", lastIDErr)
	}

	entry.SetID(uint64(id))
	return uint64(id), nil
}

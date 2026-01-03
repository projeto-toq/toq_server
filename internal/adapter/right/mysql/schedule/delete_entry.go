package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteEntry removes a single agenda entry by id.
//
// Parameters:
//   - ctx: request-scoped context for tracing/logging.
//   - tx: required transaction for atomicity.
//   - entryID: target entry identifier.
//
// Returns: sql.ErrNoRows when no row matches; driver errors for exec/rows affected problems.
// Observability: tracer span, logger propagation, span error marking on infra failures.
func (a *ScheduleAdapter) DeleteEntry(ctx context.Context, tx *sql.Tx, entryID uint64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM listing_agenda_entries WHERE id = ?`
	result, execErr := a.ExecContext(ctx, tx, "delete", query, entryID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.schedule.delete_entry.exec_error", "entry_id", entryID, "err", execErr)
		return fmt.Errorf("delete agenda entry: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.schedule.delete_entry.rows_error", "entry_id", entryID, "err", rowsErr)
		return fmt.Errorf("agenda entry rows affected: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

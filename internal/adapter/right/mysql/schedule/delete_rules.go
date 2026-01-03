package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteRulesByAgenda hard-deletes all rules for an agenda; tx required to keep writes atomic.
// Returns sql.ErrNoRows when no rule was removed; other infra errors are bubbled for service mapping.
func (a *ScheduleAdapter) DeleteRulesByAgenda(ctx context.Context, tx *sql.Tx, agendaID uint64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM listing_agenda_rules WHERE agenda_id = ?`
	result, execErr := a.ExecContext(ctx, tx, "delete", query, agendaID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.schedule.delete_rules.exec_error", "agenda_id", agendaID, "err", execErr)
		return fmt.Errorf("delete agenda rules: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.schedule.delete_rules.rows_error", "agenda_id", agendaID, "err", rowsErr)
		return fmt.Errorf("agenda rules rows affected: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

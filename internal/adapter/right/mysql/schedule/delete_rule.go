package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteRule removes a single agenda rule by id; tx required to keep write atomic.
// Returns sql.ErrNoRows when no row matches; other infra errors are logged and bubbled.
func (a *ScheduleAdapter) DeleteRule(ctx context.Context, tx *sql.Tx, ruleID uint64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM listing_agenda_rules WHERE id = ?`
	result, execErr := a.ExecContext(ctx, tx, "delete", query, ruleID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.schedule.delete_rule.exec_error", "rule_id", ruleID, "err", execErr)
		return fmt.Errorf("delete agenda rule: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.schedule.delete_rule.rows_error", "rule_id", ruleID, "err", rowsErr)
		return fmt.Errorf("agenda rule rows affected: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

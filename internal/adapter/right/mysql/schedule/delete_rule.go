package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteRule removes a single agenda rule by id.
func (a *ScheduleAdapter) DeleteRule(ctx context.Context, tx *sql.Tx, ruleID uint64) error {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM listing_agenda_rules WHERE id = ?`
	if _, err = exec.ExecContext(ctx, query, ruleID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.delete_rule.exec_error", "rule_id", ruleID, "err", err)
		return fmt.Errorf("delete agenda rule: %w", err)
	}

	return nil
}

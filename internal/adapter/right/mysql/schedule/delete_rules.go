package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ScheduleAdapter) DeleteRulesByAgenda(ctx context.Context, tx *sql.Tx, agendaID uint64) error {
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

	query := `DELETE FROM listing_agenda_rules WHERE agenda_id = ?`
	if _, err = exec.ExecContext(ctx, query, agendaID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.delete_rules.exec_error", "agenda_id", agendaID, "err", err)
		return fmt.Errorf("delete agenda rules: %w", err)
	}

	return nil
}

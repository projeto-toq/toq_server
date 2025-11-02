package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ScheduleAdapter) DeleteRulesByAgenda(ctx context.Context, tx *sql.Tx, agendaID uint64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM listing_agenda_rules WHERE agenda_id = ?`
	if _, execErr := a.ExecContext(ctx, tx, "delete", query, agendaID); execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.schedule.delete_rules.exec_error", "agenda_id", agendaID, "err", execErr)
		return fmt.Errorf("delete agenda rules: %w", execErr)
	}

	return nil
}

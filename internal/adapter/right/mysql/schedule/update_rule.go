package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateRule updates an existing agenda rule definition.
func (a *ScheduleAdapter) UpdateRule(ctx context.Context, tx *sql.Tx, rule schedulemodel.AgendaRuleInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE listing_agenda_rules SET agenda_id = ?, day_of_week = ?, start_minute = ?, end_minute = ?, rule_type = ?, is_active = ? WHERE id = ?`

	if _, execErr := a.ExecContext(ctx, tx, "update", query, rule.AgendaID(), rule.DayOfWeek(), rule.StartMinutes(), rule.EndMinutes(), rule.RuleType(), rule.IsActive(), rule.ID()); execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.schedule.update_rule.exec_error", "rule_id", rule.ID(), "err", execErr)
		return fmt.Errorf("update agenda rule: %w", execErr)
	}

	return nil
}

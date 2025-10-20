package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ScheduleAdapter) InsertRules(ctx context.Context, tx *sql.Tx, rules []schedulemodel.AgendaRuleInterface) error {
	if len(rules) == 0 {
		return nil
	}

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

	entities := converters.ToRuleEntities(rules)

	query := `INSERT INTO listing_agenda_rules (agenda_id, day_of_week, start_minute, end_minute, rule_type, is_active) VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := exec.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.insert_rules.prepare_error", "err", err)
		return fmt.Errorf("prepare insert agenda rules: %w", err)
	}
	defer stmt.Close()

	for _, rule := range entities {
		if _, err = stmt.ExecContext(ctx, rule.AgendaID, rule.DayOfWeek, rule.StartMin, rule.EndMin, rule.RuleType, rule.IsActive); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.schedule.insert_rules.exec_error", "agenda_id", rule.AgendaID, "err", err)
			return fmt.Errorf("exec insert agenda rule: %w", err)
		}
	}

	return nil
}

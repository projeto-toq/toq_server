package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entity"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ScheduleAdapter) ListRulesByAgenda(ctx context.Context, tx *sql.Tx, agendaID uint64) ([]schedulemodel.AgendaRuleInterface, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return nil, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, agenda_id, day_of_week, start_minute, end_minute, rule_type, is_active FROM listing_agenda_rules WHERE agenda_id = ? ORDER BY day_of_week, start_minute`

	rows, err := exec.QueryContext(ctx, query, agendaID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.list_rules.query_error", "agenda_id", agendaID, "err", err)
		return nil, fmt.Errorf("query agenda rules: %w", err)
	}
	defer rows.Close()

	var results []schedulemodel.AgendaRuleInterface
	for rows.Next() {
		var ruleEntity entity.RuleEntity
		if err = rows.Scan(&ruleEntity.ID, &ruleEntity.AgendaID, &ruleEntity.DayOfWeek, &ruleEntity.StartMinute, &ruleEntity.EndMinute, &ruleEntity.RuleType, &ruleEntity.IsActive); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.schedule.list_rules.scan_error", "agenda_id", agendaID, "err", err)
			return nil, fmt.Errorf("scan agenda rule: %w", err)
		}
		results = append(results, converters.ToRuleModel(ruleEntity))
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.list_rules.rows_error", "agenda_id", agendaID, "err", err)
		return nil, fmt.Errorf("iterate agenda rules: %w", err)
	}

	return results, nil
}

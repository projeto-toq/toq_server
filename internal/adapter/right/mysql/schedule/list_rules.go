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
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, agenda_id, day_of_week, start_minute, end_minute, rule_type, is_active FROM listing_agenda_rules WHERE agenda_id = ? ORDER BY day_of_week, start_minute`

	rows, queryErr := a.QueryContext(ctx, tx, "select", query, agendaID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.schedule.list_rules.query_error", "agenda_id", agendaID, "err", queryErr)
		return nil, fmt.Errorf("query agenda rules: %w", queryErr)
	}
	defer rows.Close()

	var results []schedulemodel.AgendaRuleInterface
	for rows.Next() {
		var ruleEntity entity.RuleEntity
		if scanErr := rows.Scan(&ruleEntity.ID, &ruleEntity.AgendaID, &ruleEntity.DayOfWeek, &ruleEntity.StartMinute, &ruleEntity.EndMinute, &ruleEntity.RuleType, &ruleEntity.IsActive); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.schedule.list_rules.scan_error", "agenda_id", agendaID, "err", scanErr)
			return nil, fmt.Errorf("scan agenda rule: %w", scanErr)
		}
		results = append(results, converters.ToRuleModel(ruleEntity))
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.schedule.list_rules.rows_error", "agenda_id", agendaID, "err", rowsErr)
		return nil, fmt.Errorf("iterate agenda rules: %w", rowsErr)
	}

	return results, nil
}

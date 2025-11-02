package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entity"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetRuleByID retrieves a single rule by its identifier.
func (a *ScheduleAdapter) GetRuleByID(ctx context.Context, tx *sql.Tx, ruleID uint64) (schedulemodel.AgendaRuleInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, agenda_id, day_of_week, start_minute, end_minute, rule_type, is_active FROM listing_agenda_rules WHERE id = ?`

	var ruleEntity entity.RuleEntity
	if err = a.QueryRowContext(ctx, tx, "select", query, ruleID).Scan(&ruleEntity.ID, &ruleEntity.AgendaID, &ruleEntity.DayOfWeek, &ruleEntity.StartMinute, &ruleEntity.EndMinute, &ruleEntity.RuleType, &ruleEntity.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.get_rule.scan_error", "rule_id", ruleID, "err", err)
		return nil, fmt.Errorf("scan agenda rule: %w", err)
	}

	return converters.ToRuleModel(ruleEntity), nil
}

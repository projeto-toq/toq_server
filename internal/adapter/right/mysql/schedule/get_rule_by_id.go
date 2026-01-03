package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	scheduleconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	scheduleentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entities"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetRuleByID retrieves a single agenda rule by id; tx required.
// Returns sql.ErrNoRows when not found; infra errors are logged and propagated unchanged.
func (a *ScheduleAdapter) GetRuleByID(ctx context.Context, tx *sql.Tx, ruleID uint64) (schedulemodel.AgendaRuleInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, agenda_id, day_of_week, start_minute, end_minute, rule_type, is_active FROM listing_agenda_rules WHERE id = ?`

	var ruleEntity scheduleentity.RuleEntity
	if err = a.QueryRowContext(ctx, tx, "select", query, ruleID).Scan(&ruleEntity.ID, &ruleEntity.AgendaID, &ruleEntity.DayOfWeek, &ruleEntity.StartMinute, &ruleEntity.EndMinute, &ruleEntity.RuleType, &ruleEntity.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.get_rule.scan_error", "rule_id", ruleID, "err", err)
		return nil, fmt.Errorf("scan agenda rule: %w", err)
	}

	return scheduleconverters.RuleEntityToDomain(ruleEntity), nil
}

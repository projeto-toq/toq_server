package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entity"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ScheduleAdapter) ListBlockRules(ctx context.Context, tx *sql.Tx, filter schedulemodel.BlockRulesFilter) ([]schedulemodel.AgendaRuleInterface, error) {
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

	conditions := []string{"a.owner_id = ?", "a.listing_id = ?", "r.rule_type = ?"}
	args := []any{filter.OwnerID, filter.ListingID, schedulemodel.RuleTypeBlock}

	if len(filter.Weekdays) > 0 {
		placeholders, weekdayArgs := buildWeekdayConditions(filter.Weekdays)
		conditions = append(conditions, fmt.Sprintf("r.day_of_week IN (%s)", placeholders))
		args = append(args, weekdayArgs...)
	}

	query := fmt.Sprintf(`
        SELECT r.id,
               r.agenda_id,
               r.day_of_week,
               r.start_minute,
               r.end_minute,
               r.rule_type,
               r.is_active
        FROM listing_agenda_rules r
        INNER JOIN listing_agendas a ON a.id = r.agenda_id
        WHERE %s
        ORDER BY r.day_of_week, r.start_minute
    `, strings.Join(conditions, " AND "))

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.list_block_rules.query_error", "listing_id", filter.ListingID, "err", err)
		return nil, fmt.Errorf("query block rules: %w", err)
	}
	defer rows.Close()

	rules := make([]schedulemodel.AgendaRuleInterface, 0)
	for rows.Next() {
		var ruleEntity entity.RuleEntity
		if err = rows.Scan(&ruleEntity.ID, &ruleEntity.AgendaID, &ruleEntity.DayOfWeek, &ruleEntity.StartMinute, &ruleEntity.EndMinute, &ruleEntity.RuleType, &ruleEntity.IsActive); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.schedule.list_block_rules.scan_error", "listing_id", filter.ListingID, "err", err)
			return nil, fmt.Errorf("scan block rule: %w", err)
		}
		rules = append(rules, converters.ToRuleModel(ruleEntity))
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.list_block_rules.rows_error", "listing_id", filter.ListingID, "err", err)
		return nil, fmt.Errorf("iterate block rules: %w", err)
	}

	return rules, nil
}

func buildWeekdayConditions(values []time.Weekday) (string, []any) {
	if len(values) == 0 {
		return "", nil
	}
	placeholders := make([]string, 0, len(values))
	args := make([]any, 0, len(values))
	for _, weekday := range values {
		placeholders = append(placeholders, "?")
		args = append(args, uint8(weekday))
	}
	return strings.Join(placeholders, ","), args
}

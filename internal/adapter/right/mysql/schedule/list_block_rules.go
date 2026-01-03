package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	scheduleconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	scheduleentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entities"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListBlockRules lists blocking rules filtered by owner, listing, and optional weekdays.
//
// Parameters:
//   - ctx: request-scoped context for tracing/logging.
//   - tx: optional transaction for consistent reads.
//   - filter: owner/listing identifiers and optional weekdays slice.
//
// Returns: slice of AgendaRuleInterface (empty when none) or infrastructure errors; sql.ErrNoRows is not used.
// Observability: tracer span, logger propagation, span error marking on infra failures.
func (a *ScheduleAdapter) ListBlockRules(ctx context.Context, tx *sql.Tx, filter schedulemodel.BlockRulesFilter) ([]schedulemodel.AgendaRuleInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	conditions := []string{"a.owner_id = ?", "a.listing_identity_id = ?", "r.rule_type = ?"}
	args := []any{filter.OwnerID, filter.ListingIdentityID, schedulemodel.RuleTypeBlock}

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

	rows, queryErr := a.QueryContext(ctx, tx, "select", query, args...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.schedule.list_block_rules.query_error", "listing_identity_id", filter.ListingIdentityID, "err", queryErr)
		return nil, fmt.Errorf("query block rules: %w", queryErr)
	}
	defer rows.Close()

	rules := make([]schedulemodel.AgendaRuleInterface, 0)
	for rows.Next() {
		var ruleEntity scheduleentity.RuleEntity
		if scanErr := rows.Scan(&ruleEntity.ID, &ruleEntity.AgendaID, &ruleEntity.DayOfWeek, &ruleEntity.StartMinute, &ruleEntity.EndMinute, &ruleEntity.RuleType, &ruleEntity.IsActive); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.schedule.list_block_rules.scan_error", "listing_identity_id", filter.ListingIdentityID, "err", scanErr)
			return nil, fmt.Errorf("scan block rule: %w", scanErr)
		}
		rules = append(rules, scheduleconverters.RuleEntityToDomain(ruleEntity))
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.schedule.list_block_rules.rows_error", "listing_identity_id", filter.ListingIdentityID, "err", rowsErr)
		return nil, fmt.Errorf("iterate block rules: %w", rowsErr)
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

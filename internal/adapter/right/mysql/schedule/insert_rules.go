package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entity"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ScheduleAdapter) InsertRules(ctx context.Context, tx *sql.Tx, rules []schedulemodel.AgendaRuleInterface) error {
	if len(rules) == 0 {
		return nil
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO listing_agenda_rules (agenda_id, day_of_week, start_minute, end_minute, rule_type, is_active) VALUES (?, ?, ?, ?, ?, ?)`
	stmt, cleanup, prepareErr := a.PrepareContext(ctx, tx, "insert", query)
	if prepareErr != nil {
		utils.SetSpanError(ctx, prepareErr)
		logger.Error("mysql.schedule.insert_rules.prepare_error", "err", prepareErr)
		return fmt.Errorf("prepare insert agenda rules: %w", prepareErr)
	}
	defer cleanup()

	for _, rule := range rules {
		record := entity.RuleEntity{
			AgendaID:    rule.AgendaID(),
			DayOfWeek:   uint8(rule.DayOfWeek()),
			StartMinute: rule.StartMinutes(),
			EndMinute:   rule.EndMinutes(),
			RuleType:    string(rule.RuleType()),
			IsActive:    rule.IsActive(),
		}
		result, execErr := stmt.ExecContext(ctx, record.AgendaID, record.DayOfWeek, record.StartMinute, record.EndMinute, record.RuleType, record.IsActive)
		if execErr != nil {
			utils.SetSpanError(ctx, execErr)
			logger.Error("mysql.schedule.insert_rules.exec_error", "agenda_id", record.AgendaID, "err", execErr)
			return fmt.Errorf("exec insert agenda rule: %w", execErr)
		}
		if lastID, idErr := result.LastInsertId(); idErr == nil {
			rule.SetID(uint64(lastID))
		} else {
			utils.SetSpanError(ctx, idErr)
			logger.Error("mysql.schedule.insert_rules.last_id_error", "agenda_id", record.AgendaID, "err", idErr)
			return fmt.Errorf("retrieve agenda rule id: %w", idErr)
		}
	}

	return nil
}

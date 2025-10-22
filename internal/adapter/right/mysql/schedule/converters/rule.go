package converters

import (
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entity"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
)

// ToRuleModel converts RuleEntity to domain.
func ToRuleModel(e entity.RuleEntity) schedulemodel.AgendaRuleInterface {
	rule := schedulemodel.NewAgendaRule()
	rule.SetID(e.ID)
	rule.SetAgendaID(e.AgendaID)
	rule.SetDayOfWeek(time.Weekday(e.DayOfWeek))
	rule.SetStartMinutes(e.StartMinute)
	rule.SetEndMinutes(e.EndMinute)
	rule.SetRuleType(schedulemodel.RuleType(e.RuleType))
	rule.SetActive(e.IsActive)
	return rule
}

// ToRuleEntities converts domain models to persistence.
func ToRuleEntities(models []schedulemodel.AgendaRuleInterface) []entity.RuleEntity {
	entities := make([]entity.RuleEntity, 0, len(models))
	for _, model := range models {
		entities = append(entities, entity.RuleEntity{
			ID:          model.ID(),
			AgendaID:    model.AgendaID(),
			DayOfWeek:   uint8(model.DayOfWeek()),
			StartMinute: model.StartMinutes(),
			EndMinute:   model.EndMinutes(),
			RuleType:    string(model.RuleType()),
			IsActive:    model.IsActive(),
		})
	}
	return entities
}

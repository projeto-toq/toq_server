package scheduleconverters

import (
	"time"

	scheduleentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entities"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
)

// RuleEntityToDomain converts RuleEntity to the domain representation.
func RuleEntityToDomain(e scheduleentity.RuleEntity) schedulemodel.AgendaRuleInterface {
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

// RuleDomainsToEntities converts domain rules into persistence entities.
func RuleDomainsToEntities(models []schedulemodel.AgendaRuleInterface) []scheduleentity.RuleEntity {
	entities := make([]scheduleentity.RuleEntity, 0, len(models))
	for _, model := range models {
		entities = append(entities, scheduleentity.RuleEntity{
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

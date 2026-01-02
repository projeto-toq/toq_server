package scheduleservices

import (
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
)

// normalizeAvailabilityRulesToExclusive converts stored inclusive end minutes into half-open ranges.
// listing_agenda_rules persists end_minute as the last blocked minute (inclusive). The availability engine
// works with half-open intervals [start, end), so we bump the end minute by 1 (clamped at 1440) for
// blocking rules before computing availability.
func normalizeAvailabilityRulesToExclusive(rules []schedulemodel.AgendaRuleInterface) {
    for _, rule := range rules {
        if rule == nil {
            continue
        }
        if rule.RuleType() != schedulemodel.RuleTypeBlock {
            continue
        }

        end := rule.EndMinutes()
        if end < minutesPerDay {
            end++
        } else {
            end = minutesPerDay
        }
        rule.SetEndMinutes(end)
    }
}

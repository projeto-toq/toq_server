package schedulemodel

import "time"

// AgendaRuleInterface represents a recurring rule applied to an agenda.
type AgendaRuleInterface interface {
	ID() uint64
	SetID(id uint64)
	AgendaID() uint64
	SetAgendaID(agendaID uint64)
	DayOfWeek() time.Weekday
	SetDayOfWeek(value time.Weekday)
	StartMinutes() uint16
	SetStartMinutes(value uint16)
	EndMinutes() uint16
	SetEndMinutes(value uint16)
	RuleType() RuleType
	SetRuleType(value RuleType)
	IsActive() bool
	SetActive(value bool)
}

// NewAgendaRule returns a new rule object.
func NewAgendaRule() AgendaRuleInterface {
	return &agendaRule{}
}

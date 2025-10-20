package schedulemodel

import "time"

type agendaRule struct {
	id          uint64
	agendaID    uint64
	dayOfWeek   time.Weekday
	startMinute uint16
	endMinute   uint16
	ruleType    RuleType
	active      bool
}

func (r *agendaRule) ID() uint64 {
	return r.id
}

func (r *agendaRule) SetID(id uint64) {
	r.id = id
}

func (r *agendaRule) AgendaID() uint64 {
	return r.agendaID
}

func (r *agendaRule) SetAgendaID(agendaID uint64) {
	r.agendaID = agendaID
}

func (r *agendaRule) DayOfWeek() time.Weekday {
	return r.dayOfWeek
}

func (r *agendaRule) SetDayOfWeek(value time.Weekday) {
	r.dayOfWeek = value
}

func (r *agendaRule) StartMinutes() uint16 {
	return r.startMinute
}

func (r *agendaRule) SetStartMinutes(value uint16) {
	r.startMinute = value
}

func (r *agendaRule) EndMinutes() uint16 {
	return r.endMinute
}

func (r *agendaRule) SetEndMinutes(value uint16) {
	r.endMinute = value
}

func (r *agendaRule) RuleType() RuleType {
	return r.ruleType
}

func (r *agendaRule) SetRuleType(value RuleType) {
	r.ruleType = value
}

func (r *agendaRule) IsActive() bool {
	return r.active
}

func (r *agendaRule) SetActive(value bool) {
	r.active = value
}

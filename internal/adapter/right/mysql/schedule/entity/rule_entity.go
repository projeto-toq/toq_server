package entity

// RuleEntity mirrors the listing_agenda_rules table.
type RuleEntity struct {
	ID          uint64
	AgendaID    uint64
	DayOfWeek   uint8
	StartMinute uint16
	EndMinute   uint16
	RuleType    string
	IsActive    bool
}

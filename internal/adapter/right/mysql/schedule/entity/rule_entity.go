package entity

// RuleEntity mirrors the listing_agenda_rules table.
type RuleEntity struct {
	ID        uint64
	AgendaID  uint64
	DayOfWeek uint8
	StartMin  uint16
	EndMin    uint16
	RuleType  string
	IsActive  bool
}

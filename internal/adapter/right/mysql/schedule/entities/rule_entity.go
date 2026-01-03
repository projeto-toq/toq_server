package scheduleentity

// RuleEntity maps listing_agenda_rules rows (unique per agenda/day/start/end) for adapter use only.
type RuleEntity struct {
	ID          uint64 // PK AUTO_INCREMENT
	AgendaID    uint64 // FK listing_agendas.id (INT UNSIGNED NOT NULL)
	DayOfWeek   uint8  // TINYINT NOT NULL, 0=Sunday..6=Saturday
	StartMinute uint16 // INT UNSIGNED NOT NULL, minutes since 00:00
	EndMinute   uint16 // INT UNSIGNED NOT NULL, minutes since 00:00
	RuleType    string // ENUM('BLOCK','FREE') NOT NULL
	IsActive    bool   // TINYINT(1) NOT NULL DEFAULT 1
}

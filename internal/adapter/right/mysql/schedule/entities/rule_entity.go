package scheduleentity

// RuleEntity maps listing_agenda_rules rows (unique per agenda/day/start/end) for adapter use only.
// Schema summary: PK (id), FK agenda_id -> listing_agendas.id, unique (agenda_id, day_of_week, start_minute, end_minute), default is_active=1.
// Conversions to/from domain are handled by scheduleconverters; no domain imports should be added here.
type RuleEntity struct {
	// ID is the AUTO_INCREMENT primary key (INT UNSIGNED NOT NULL).
	ID uint64
	// AgendaID references listing_agendas.id (INT UNSIGNED NOT NULL, indexed fk_rules_agenda_idx).
	AgendaID uint64
	// DayOfWeek is the weekday number (TINYINT NOT NULL, 0=Sunday..6=Saturday).
	DayOfWeek uint8
	// StartMinute is the rule start offset in minutes from 00:00 (INT UNSIGNED NOT NULL).
	StartMinute uint16
	// EndMinute is the rule end offset in minutes from 00:00 (INT UNSIGNED NOT NULL).
	EndMinute uint16
	// RuleType stores the rule effect (ENUM('BLOCK','FREE') NOT NULL).
	RuleType string
	// IsActive indicates whether the rule is currently applied (TINYINT(1) NOT NULL DEFAULT 1).
	IsActive bool
}

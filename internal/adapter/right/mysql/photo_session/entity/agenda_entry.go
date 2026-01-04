package entity

import "database/sql"

// AgendaEntry models photographer_agenda_entries with nullable columns explicitly represented.
// Columns: id (PK, NOT NULL), photographer_user_id (NOT NULL), entry_type (NOT NULL),
// source (ENUM NULL), source_id (INT NULL), starts_at (DATETIME(6) NULL), ends_at (DATETIME(6) NULL),
// blocking (TINYINT NULL DEFAULT 1), reason (VARCHAR(255) NULL), timezone (VARCHAR(50) NULL DEFAULT 'America/Sao_Paulo').
type AgendaEntry struct {
	ID                 uint64         // photographer_agenda_entries.id
	PhotographerUserID uint64         // photographer_agenda_entries.photographer_user_id (NOT NULL)
	EntryType          string         // photographer_agenda_entries.entry_type (ENUM, NOT NULL)
	Source             sql.NullString // photographer_agenda_entries.source (ENUM, NULL)
	SourceID           sql.NullInt64  // photographer_agenda_entries.source_id (NULLABLE)
	StartsAt           sql.NullTime   // photographer_agenda_entries.starts_at (DATETIME(6), NULLABLE)
	EndsAt             sql.NullTime   // photographer_agenda_entries.ends_at (DATETIME(6), NULLABLE)
	Blocking           sql.NullBool   // photographer_agenda_entries.blocking (TINYINT, NULL DEFAULT 1)
	Reason             sql.NullString // photographer_agenda_entries.reason (VARCHAR(255), NULLABLE)
	Timezone           sql.NullString // photographer_agenda_entries.timezone (VARCHAR(50), NULL DEFAULT 'America/Sao_Paulo')
}

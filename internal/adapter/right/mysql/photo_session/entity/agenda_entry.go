package entity

import (
	"database/sql"
	"time"
)

// AgendaEntry represents a row from photographer_agenda_entries.
type AgendaEntry struct {
	ID                 uint64
	PhotographerUserID uint64
	EntryType          string
	Source             string
	SourceID           sql.NullInt64
	StartsAt           time.Time
	EndsAt             time.Time
	Blocking           bool
	Reason             sql.NullString
	Timezone           string
	CreatedAt          sql.NullTime
	UpdatedAt          sql.NullTime
}

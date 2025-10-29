package photosessionmodel

import "time"

// AgendaEntryInterface describes a single item stored in photographer_agenda_entries.
type AgendaEntryInterface interface {
	ID() uint64
	SetID(id uint64)
	PhotographerUserID() uint64
	SetPhotographerUserID(id uint64)
	EntryType() AgendaEntryType
	SetEntryType(t AgendaEntryType)
	Source() AgendaEntrySource
	SetSource(source AgendaEntrySource)
	SourceID() (*uint64, bool)
	SetSourceID(id uint64)
	ClearSourceID()
	StartsAt() time.Time
	SetStartsAt(t time.Time)
	EndsAt() time.Time
	SetEndsAt(t time.Time)
	Blocking() bool
	SetBlocking(blocking bool)
	Reason() (string, bool)
	SetReason(reason string)
	ClearReason()
	Timezone() string
	SetTimezone(tz string)
	CreatedAt() (time.Time, bool)
	SetCreatedAt(t time.Time)
	UpdatedAt() (time.Time, bool)
	SetUpdatedAt(t time.Time)
}

// NewAgendaEntry builds a new mutable agenda entry.
func NewAgendaEntry() AgendaEntryInterface {
	return &agendaEntry{}
}

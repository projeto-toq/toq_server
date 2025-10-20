package schedulemodel

import "time"

// AgendaEntryInterface represents a concrete entry in an agenda.
type AgendaEntryInterface interface {
	ID() uint64
	SetID(id uint64)
	AgendaID() uint64
	SetAgendaID(agendaID uint64)
	EntryType() EntryType
	SetEntryType(value EntryType)
	StartsAt() time.Time
	SetStartsAt(value time.Time)
	EndsAt() time.Time
	SetEndsAt(value time.Time)
	Blocking() bool
	SetBlocking(value bool)
	Reason() (string, bool)
	SetReason(value string)
	ClearReason()
	VisitID() (uint64, bool)
	SetVisitID(value uint64)
	ClearVisitID()
	PhotoBookingID() (uint64, bool)
	SetPhotoBookingID(value uint64)
	ClearPhotoBookingID()
	CreatedBy() int64
	SetCreatedBy(value int64)
	UpdatedBy() (int64, bool)
	SetUpdatedBy(value int64)
	ClearUpdatedBy()
}

// NewAgendaEntry builds a new agenda entry instance.
func NewAgendaEntry() AgendaEntryInterface {
	return &agendaEntry{}
}

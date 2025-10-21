package schedulemodel

import "time"

type agendaEntry struct {
	id             uint64
	agendaID       uint64
	entryType      EntryType
	startsAt       time.Time
	endsAt         time.Time
	blocking       bool
	reason         string
	reasonValid    bool
	visitID        uint64
	visitValid     bool
	photoBookingID uint64
	photoValid     bool
}

func (e *agendaEntry) ID() uint64 {
	return e.id
}

func (e *agendaEntry) SetID(id uint64) {
	e.id = id
}

func (e *agendaEntry) AgendaID() uint64 {
	return e.agendaID
}

func (e *agendaEntry) SetAgendaID(agendaID uint64) {
	e.agendaID = agendaID
}

func (e *agendaEntry) EntryType() EntryType {
	return e.entryType
}

func (e *agendaEntry) SetEntryType(value EntryType) {
	e.entryType = value
}

func (e *agendaEntry) StartsAt() time.Time {
	return e.startsAt
}

func (e *agendaEntry) SetStartsAt(value time.Time) {
	e.startsAt = value
}

func (e *agendaEntry) EndsAt() time.Time {
	return e.endsAt
}

func (e *agendaEntry) SetEndsAt(value time.Time) {
	e.endsAt = value
}

func (e *agendaEntry) Blocking() bool {
	return e.blocking
}

func (e *agendaEntry) SetBlocking(value bool) {
	e.blocking = value
}

func (e *agendaEntry) Reason() (string, bool) {
	return e.reason, e.reasonValid
}

func (e *agendaEntry) SetReason(value string) {
	e.reason = value
	e.reasonValid = true
}

func (e *agendaEntry) ClearReason() {
	e.reason = ""
	e.reasonValid = false
}

func (e *agendaEntry) VisitID() (uint64, bool) {
	return e.visitID, e.visitValid
}

func (e *agendaEntry) SetVisitID(value uint64) {
	e.visitID = value
	e.visitValid = true
}

func (e *agendaEntry) ClearVisitID() {
	e.visitID = 0
	e.visitValid = false
}

func (e *agendaEntry) PhotoBookingID() (uint64, bool) {
	return e.photoBookingID, e.photoValid
}

func (e *agendaEntry) SetPhotoBookingID(value uint64) {
	e.photoBookingID = value
	e.photoValid = true
}

func (e *agendaEntry) ClearPhotoBookingID() {
	e.photoBookingID = 0
	e.photoValid = false
}

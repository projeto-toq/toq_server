package photosessionmodel

import "time"

type agendaEntry struct {
	id                 uint64
	photographerUserID uint64
	entryType          AgendaEntryType
	source             AgendaEntrySource
	sourceID           uint64
	sourceIDValid      bool
	startsAt           time.Time
	endsAt             time.Time
	blocking           bool
	reason             string
	reasonValid        bool
	timezone           string
	createdAt          time.Time
	createdAtValid     bool
	updatedAt          time.Time
	updatedAtValid     bool
}

func (e *agendaEntry) ID() uint64 {
	return e.id
}

func (e *agendaEntry) SetID(id uint64) {
	e.id = id
}

func (e *agendaEntry) PhotographerUserID() uint64 {
	return e.photographerUserID
}

func (e *agendaEntry) SetPhotographerUserID(id uint64) {
	e.photographerUserID = id
}

func (e *agendaEntry) EntryType() AgendaEntryType {
	return e.entryType
}

func (e *agendaEntry) SetEntryType(t AgendaEntryType) {
	e.entryType = t
}

func (e *agendaEntry) Source() AgendaEntrySource {
	return e.source
}

func (e *agendaEntry) SetSource(source AgendaEntrySource) {
	e.source = source
}

func (e *agendaEntry) SourceID() (*uint64, bool) {
	if !e.sourceIDValid {
		return nil, false
	}
	return &e.sourceID, true
}

func (e *agendaEntry) SetSourceID(id uint64) {
	e.sourceID = id
	e.sourceIDValid = true
}

func (e *agendaEntry) ClearSourceID() {
	e.sourceID = 0
	e.sourceIDValid = false
}

func (e *agendaEntry) StartsAt() time.Time {
	return e.startsAt
}

func (e *agendaEntry) SetStartsAt(t time.Time) {
	e.startsAt = t
}

func (e *agendaEntry) EndsAt() time.Time {
	return e.endsAt
}

func (e *agendaEntry) SetEndsAt(t time.Time) {
	e.endsAt = t
}

func (e *agendaEntry) Blocking() bool {
	return e.blocking
}

func (e *agendaEntry) SetBlocking(blocking bool) {
	e.blocking = blocking
}

func (e *agendaEntry) Reason() (string, bool) {
	if !e.reasonValid {
		return "", false
	}
	return e.reason, true
}

func (e *agendaEntry) SetReason(reason string) {
	e.reason = reason
	e.reasonValid = true
}

func (e *agendaEntry) ClearReason() {
	e.reason = ""
	e.reasonValid = false
}

func (e *agendaEntry) Timezone() string {
	return e.timezone
}

func (e *agendaEntry) SetTimezone(tz string) {
	e.timezone = tz
}

func (e *agendaEntry) CreatedAt() (time.Time, bool) {
	if !e.createdAtValid {
		return time.Time{}, false
	}
	return e.createdAt, true
}

func (e *agendaEntry) SetCreatedAt(t time.Time) {
	e.createdAt = t
	e.createdAtValid = true
}

func (e *agendaEntry) UpdatedAt() (time.Time, bool) {
	if !e.updatedAtValid {
		return time.Time{}, false
	}
	return e.updatedAt, true
}

func (e *agendaEntry) SetUpdatedAt(t time.Time) {
	e.updatedAt = t
	e.updatedAtValid = true
}

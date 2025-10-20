package listingmodel

import "time"

type visit struct {
	id             int64
	listingID      int64
	ownerID        int64
	realtorID      int64
	scheduledStart time.Time
	scheduledEnd   time.Time
	status         VisitStatus
	cancelReason   string
	cancelValid    bool
	notes          string
	notesValid     bool
	createdBy      int64
	updatedBy      int64
	updatedValid   bool
}

func (v *visit) ID() int64 {
	return v.id
}

func (v *visit) SetID(id int64) {
	v.id = id
}

func (v *visit) ListingID() int64 {
	return v.listingID
}

func (v *visit) SetListingID(id int64) {
	v.listingID = id
}

func (v *visit) OwnerID() int64 {
	return v.ownerID
}

func (v *visit) SetOwnerID(id int64) {
	v.ownerID = id
}

func (v *visit) RealtorID() int64 {
	return v.realtorID
}

func (v *visit) SetRealtorID(id int64) {
	v.realtorID = id
}

func (v *visit) ScheduledStart() time.Time {
	return v.scheduledStart
}

func (v *visit) SetScheduledStart(value time.Time) {
	v.scheduledStart = value
}

func (v *visit) ScheduledEnd() time.Time {
	return v.scheduledEnd
}

func (v *visit) SetScheduledEnd(value time.Time) {
	v.scheduledEnd = value
}

func (v *visit) Status() VisitStatus {
	return v.status
}

func (v *visit) SetStatus(value VisitStatus) {
	v.status = value
}

func (v *visit) CancelReason() (string, bool) {
	return v.cancelReason, v.cancelValid
}

func (v *visit) SetCancelReason(value string) {
	v.cancelReason = value
	v.cancelValid = true
}

func (v *visit) ClearCancelReason() {
	v.cancelReason = ""
	v.cancelValid = false
}

func (v *visit) Notes() (string, bool) {
	return v.notes, v.notesValid
}

func (v *visit) SetNotes(value string) {
	v.notes = value
	v.notesValid = true
}

func (v *visit) ClearNotes() {
	v.notes = ""
	v.notesValid = false
}

func (v *visit) CreatedBy() int64 {
	return v.createdBy
}

func (v *visit) SetCreatedBy(value int64) {
	v.createdBy = value
}

func (v *visit) UpdatedBy() (int64, bool) {
	return v.updatedBy, v.updatedValid
}

func (v *visit) SetUpdatedBy(value int64) {
	v.updatedBy = value
	v.updatedValid = true
}

func (v *visit) ClearUpdatedBy() {
	v.updatedBy = 0
	v.updatedValid = false
}

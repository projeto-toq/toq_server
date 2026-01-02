package listingmodel

import "time"

type visit struct {
	id                int64
	listingIdentityID int64
	listingVersion    uint8
	requesterUserID   int64
	ownerUserID       int64
	scheduledStart    time.Time
	scheduledEnd      time.Time
	status            VisitStatus
	source            string
	sourceValid       bool
	notes             string
	notesValid        bool
	rejectionReason   string
	rejectionValid    bool
	firstOwnerAction  time.Time
	firstOwnerValid   bool
	requestedAt       time.Time
	createdBy         int64
	updatedBy         int64
	updatedValid      bool
	createdAt         time.Time
	updatedAt         time.Time
}

func (v *visit) ID() int64 {
	return v.id
}

func (v *visit) SetID(id int64) {
	v.id = id
}

func (v *visit) ListingIdentityID() int64 {
	return v.listingIdentityID
}

func (v *visit) SetListingIdentityID(id int64) {
	v.listingIdentityID = id
}

func (v *visit) ListingVersion() uint8 {
	return v.listingVersion
}

func (v *visit) SetListingVersion(version uint8) {
	v.listingVersion = version
}

func (v *visit) RequesterUserID() int64 {
	return v.requesterUserID
}

func (v *visit) SetRequesterUserID(id int64) {
	v.requesterUserID = id
}

// UserID is kept as an alias to RequesterUserID for compatibility.
func (v *visit) UserID() int64 {
	return v.requesterUserID
}

func (v *visit) OwnerUserID() int64 {
	return v.ownerUserID
}

func (v *visit) SetOwnerUserID(id int64) {
	v.ownerUserID = id
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

func (v *visit) Source() (string, bool) {
	return v.source, v.sourceValid
}

func (v *visit) SetSource(value string) {
	v.source = value
	v.sourceValid = true
}

func (v *visit) ClearSource() {
	v.source = ""
	v.sourceValid = false
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

func (v *visit) RejectionReason() (string, bool) {
	return v.rejectionReason, v.rejectionValid
}

func (v *visit) SetRejectionReason(value string) {
	v.rejectionReason = value
	v.rejectionValid = true
}

func (v *visit) ClearRejectionReason() {
	v.rejectionReason = ""
	v.rejectionValid = false
}

func (v *visit) FirstOwnerActionAt() (time.Time, bool) {
	return v.firstOwnerAction, v.firstOwnerValid
}

func (v *visit) SetFirstOwnerActionAt(value time.Time) {
	v.firstOwnerAction = value
	v.firstOwnerValid = true
}

func (v *visit) ClearFirstOwnerActionAt() {
	v.firstOwnerAction = time.Time{}
	v.firstOwnerValid = false
}

func (v *visit) RequestedAt() time.Time {
	return v.requestedAt
}

func (v *visit) SetRequestedAt(value time.Time) {
	v.requestedAt = value
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

func (v *visit) CreatedAt() time.Time {
	return v.createdAt
}

func (v *visit) SetCreatedAt(value time.Time) {
	v.createdAt = value
}

func (v *visit) UpdatedAt() time.Time {
	return v.updatedAt
}

func (v *visit) SetUpdatedAt(value time.Time) {
	v.updatedAt = value
}

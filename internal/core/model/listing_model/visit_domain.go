package listingmodel

import "time"

type VisitMode string

const (
	VisitModeWithClient        VisitMode = "WITH_CLIENT"
	VisitModeRealtorOnly       VisitMode = "REALTOR_ONLY"
	VisitModeContentProduction VisitMode = "CONTENT_PRODUCTION"
)

type visit struct {
	id                int64
	listingIdentityID int64
	listingVersion    uint8
	requesterUserID   int64
	ownerUserID       int64
	scheduledStart    time.Time
	scheduledEnd      time.Time
	durationMinutes   int64
	status            VisitStatus
	visitMode         VisitMode
	source            string
	sourceValid       bool
	realtorNotes      string
	realtorNotesValid bool
	ownerNotes        string
	ownerNotesValid   bool
	rejectionReason   string
	rejectionValid    bool
	cancelReason      string
	cancelValid       bool
	firstOwnerAction  time.Time
	firstOwnerValid   bool
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

func (v *visit) DurationMinutes() int64 {
	return v.durationMinutes
}

func (v *visit) SetDurationMinutes(value int64) {
	v.durationMinutes = value
}

func (v *visit) Status() VisitStatus {
	return v.status
}

func (v *visit) SetStatus(value VisitStatus) {
	v.status = value
}

func (v *visit) Type() VisitMode {
	return v.visitMode
}

func (v *visit) SetType(value VisitMode) {
	v.visitMode = value
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

func (v *visit) RealtorNotes() (string, bool) {
	return v.realtorNotes, v.realtorNotesValid
}

func (v *visit) SetRealtorNotes(value string) {
	v.realtorNotes = value
	v.realtorNotesValid = true
}

func (v *visit) ClearRealtorNotes() {
	v.realtorNotes = ""
	v.realtorNotesValid = false
}

func (v *visit) OwnerNotes() (string, bool) {
	return v.ownerNotes, v.ownerNotesValid
}

func (v *visit) SetOwnerNotes(value string) {
	v.ownerNotes = value
	v.ownerNotesValid = true
}

func (v *visit) ClearOwnerNotes() {
	v.ownerNotes = ""
	v.ownerNotesValid = false
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

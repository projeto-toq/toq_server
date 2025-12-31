package listingmodel

import "time"

// VisitInterface represents a visit domain model with all mutable fields.
// Keep this interface to preserve compatibility with existing adapters while we
// evolve the visit domain.
type VisitInterface interface {
	ID() int64
	SetID(id int64)

	ListingIdentityID() int64
	SetListingIdentityID(id int64)

	ListingVersion() uint8
	SetListingVersion(version uint8)

	RequesterUserID() int64
	SetRequesterUserID(id int64)

	// UserID is maintained for backward compatibility and returns the requester ID.
	UserID() int64

	OwnerUserID() int64
	SetOwnerUserID(id int64)

	ScheduledStart() time.Time
	SetScheduledStart(value time.Time)

	ScheduledEnd() time.Time
	SetScheduledEnd(value time.Time)

	DurationMinutes() int64
	SetDurationMinutes(value int64)

	Status() VisitStatus
	SetStatus(value VisitStatus)

	Type() VisitMode
	SetType(value VisitMode)

	Source() (string, bool)
	SetSource(value string)
	ClearSource()

	RealtorNotes() (string, bool)
	SetRealtorNotes(value string)
	ClearRealtorNotes()

	OwnerNotes() (string, bool)
	SetOwnerNotes(value string)
	ClearOwnerNotes()

	RejectionReason() (string, bool)
	SetRejectionReason(value string)
	ClearRejectionReason()

	CancelReason() (string, bool)
	SetCancelReason(value string)
	ClearCancelReason()

	FirstOwnerActionAt() (time.Time, bool)
	SetFirstOwnerActionAt(value time.Time)
	ClearFirstOwnerActionAt()

	CreatedBy() int64
	SetCreatedBy(value int64)

	UpdatedBy() (int64, bool)
	SetUpdatedBy(value int64)
	ClearUpdatedBy()

	CreatedAt() time.Time
	SetCreatedAt(value time.Time)

	UpdatedAt() time.Time
	SetUpdatedAt(value time.Time)
}

func NewVisit() VisitInterface {
	return &visit{}
}

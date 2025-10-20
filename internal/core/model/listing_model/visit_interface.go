package listingmodel

import "time"

type VisitInterface interface {
	ID() int64
	SetID(id int64)
	ListingID() int64
	SetListingID(id int64)
	OwnerID() int64
	SetOwnerID(id int64)
	RealtorID() int64
	SetRealtorID(id int64)
	ScheduledStart() time.Time
	SetScheduledStart(value time.Time)
	ScheduledEnd() time.Time
	SetScheduledEnd(value time.Time)
	Status() VisitStatus
	SetStatus(value VisitStatus)
	CancelReason() (string, bool)
	SetCancelReason(value string)
	ClearCancelReason()
	Notes() (string, bool)
	SetNotes(value string)
	ClearNotes()
	CreatedBy() int64
	SetCreatedBy(value int64)
	UpdatedBy() (int64, bool)
	SetUpdatedBy(value int64)
	ClearUpdatedBy()
}

func NewVisit() VisitInterface {
	return &visit{}
}

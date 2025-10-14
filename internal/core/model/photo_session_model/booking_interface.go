package photosessionmodel

import "time"

// PhotoSessionBookingInterface defines the contract for photo session bookings.
type PhotoSessionBookingInterface interface {
	ID() uint64
	SetID(id uint64)
	SlotID() uint64
	SetSlotID(id uint64)
	ListingID() int64
	SetListingID(id int64)
	ScheduledStart() time.Time
	SetScheduledStart(value time.Time)
	ScheduledEnd() time.Time
	SetScheduledEnd(value time.Time)
	Status() BookingStatus
	SetStatus(status BookingStatus)
	CreatedBy() int64
	SetCreatedBy(id int64)
	Notes() *string
	SetNotes(notes *string)
	CreatedAt() time.Time
	SetCreatedAt(value time.Time)
	UpdatedAt() time.Time
	SetUpdatedAt(value time.Time)
}

// NewPhotoSessionBooking creates a new mutable booking instance.
func NewPhotoSessionBooking() PhotoSessionBookingInterface {
	return &photoSessionBooking{}
}

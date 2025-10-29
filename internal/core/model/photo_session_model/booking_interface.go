package photosessionmodel

import "time"

// PhotoSessionBookingInterface defines the contract for photo session bookings.
type PhotoSessionBookingInterface interface {
	ID() uint64
	SetID(id uint64)
	AgendaEntryID() uint64
	SetAgendaEntryID(id uint64)
	PhotographerUserID() uint64
	SetPhotographerUserID(id uint64)
	ListingID() int64
	SetListingID(id int64)
	StartsAt() time.Time
	SetStartsAt(value time.Time)
	EndsAt() time.Time
	SetEndsAt(value time.Time)
	Status() BookingStatus
	SetStatus(status BookingStatus)
	Reason() *string
	SetReason(reason *string)
}

// NewPhotoSessionBooking creates a new mutable booking instance.
func NewPhotoSessionBooking() PhotoSessionBookingInterface {
	return &photoSessionBooking{}
}

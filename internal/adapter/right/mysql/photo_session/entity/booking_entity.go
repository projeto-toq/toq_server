package entity

import "time"

// BookingEntity represents the photographer_slot_bookings table.
type BookingEntity struct {
	ID             uint64
	SlotID         uint64
	ListingID      int64
	ScheduledStart time.Time
	ScheduledEnd   time.Time
	Status         string
	Notes          *string
}

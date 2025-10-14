package entity

import "time"

// SlotEntity represents the photographer_time_slots table.
type SlotEntity struct {
	ID                 uint64
	PhotographerUserID uint64
	SlotDate           time.Time
	Period             string
	Status             string
	ReservationToken   *string
	ReservedUntil      *time.Time
	BookedAt           *time.Time
}

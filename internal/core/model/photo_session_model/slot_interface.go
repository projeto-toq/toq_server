package photosessionmodel

import "time"

// PhotographerSlotInterface defines a read/write contract for photographer time slots.
type PhotographerSlotInterface interface {
	ID() uint64
	SetID(id uint64)
	PhotographerUserID() uint64
	SetPhotographerUserID(id uint64)
	SlotDate() time.Time
	SetSlotDate(date time.Time)
	Period() SlotPeriod
	SetPeriod(period SlotPeriod)
	Status() SlotStatus
	SetStatus(status SlotStatus)
	ReservationToken() *string
	SetReservationToken(token *string)
	ReservedUntil() *time.Time
	SetReservedUntil(value *time.Time)
	BookedAt() *time.Time
	SetBookedAt(value *time.Time)
}

// NewPhotographerSlot creates a new mutable slot instance.
func NewPhotographerSlot() PhotographerSlotInterface {
	return &photographerSlot{}
}

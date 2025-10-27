package photosessionmodel

import "time"

// PhotographerDefaultAvailabilityInterface defines recurring availability slots for photographers.
type PhotographerDefaultAvailabilityInterface interface {
	ID() uint64
	SetID(id uint64)
	PhotographerUserID() uint64
	SetPhotographerUserID(id uint64)
	Weekday() time.Weekday
	SetWeekday(value time.Weekday)
	Period() SlotPeriod
	SetPeriod(value SlotPeriod)
	StartHour() int
	SetStartHour(value int)
	SlotsPerPeriod() int
	SetSlotsPerPeriod(value int)
	SlotDurationMinutes() int
	SetSlotDurationMinutes(value int)
}

// NewPhotographerDefaultAvailability creates a new mutable default availability record.
func NewPhotographerDefaultAvailability() PhotographerDefaultAvailabilityInterface {
	return &photographerDefaultAvailability{}
}

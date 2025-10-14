package photosessionmodel

// SlotStatus represents the current state of a photographer time slot.
type SlotStatus string

const (
	SlotStatusAvailable SlotStatus = "AVAILABLE"
	SlotStatusReserved  SlotStatus = "RESERVED"
	SlotStatusBooked    SlotStatus = "BOOKED"
	SlotStatusBlocked   SlotStatus = "BLOCKED"
)

// SlotPeriod represents the predefined period windows for photo sessions.
type SlotPeriod string

const (
	SlotPeriodMorning   SlotPeriod = "MORNING"
	SlotPeriodAfternoon SlotPeriod = "AFTERNOON"
)

// BookingStatus represents the lifecycle of a scheduled photo session.
type BookingStatus string

const (
	BookingStatusActive      BookingStatus = "ACTIVE"
	BookingStatusRescheduled BookingStatus = "RESCHEDULED"
	BookingStatusCancelled   BookingStatus = "CANCELLED"
	BookingStatusDone        BookingStatus = "DONE"
)

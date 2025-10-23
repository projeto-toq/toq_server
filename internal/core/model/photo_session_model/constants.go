package photosessionmodel

import "fmt"

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
	BookingStatusPendingApproval BookingStatus = "PENDING_APPROVAL"
	BookingStatusAccepted        BookingStatus = "ACCEPTED"
	BookingStatusRejected        BookingStatus = "REJECTED"
	BookingStatusActive          BookingStatus = "ACTIVE"
	BookingStatusRescheduled     BookingStatus = "RESCHEDULED"
	BookingStatusCancelled       BookingStatus = "CANCELLED"
	BookingStatusDone            BookingStatus = "DONE"
)

var validBookingStatus = map[BookingStatus]struct{}{
	BookingStatusPendingApproval: {},
	BookingStatusAccepted:        {},
	BookingStatusRejected:        {},
	BookingStatusActive:          {},
	BookingStatusRescheduled:     {},
	BookingStatusCancelled:       {},
	BookingStatusDone:            {},
}

// BookingStatusFromString converts a string to a BookingStatus type, returning an error if invalid.
func BookingStatusFromString(s string) (BookingStatus, error) {
	status := BookingStatus(s)
	if _, ok := validBookingStatus[status]; !ok {
		return "", fmt.Errorf("invalid booking status: %s", s)
	}
	return status, nil
}

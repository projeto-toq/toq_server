package photosessionmodel

import "fmt"

// AgendaEntryType represents the semantic meaning of an agenda entry row.
type AgendaEntryType string

const (
	AgendaEntryTypePhotoSession AgendaEntryType = "PHOTO_SESSION"
	AgendaEntryTypeBlock        AgendaEntryType = "BLOCK"
	AgendaEntryTypeTimeOff      AgendaEntryType = "TIME_OFF"
	AgendaEntryTypeHoliday      AgendaEntryType = "HOLIDAY"
)

// AgendaEntrySource tracks which pipeline originated a given entry.
type AgendaEntrySource string

const (
	AgendaEntrySourceBooking    AgendaEntrySource = "BOOKING"
	AgendaEntrySourceManual     AgendaEntrySource = "MANUAL"
	AgendaEntrySourceOnboarding AgendaEntrySource = "ONBOARDING"
	AgendaEntrySourceHoliday    AgendaEntrySource = "HOLIDAY_SYNC"
)

// SlotStatus captures the availability state of an agenda slot exposed to clients.
type SlotStatus string

const (
	SlotStatusAvailable SlotStatus = "AVAILABLE"
	SlotStatusReserved  SlotStatus = "RESERVED"
	SlotStatusBooked    SlotStatus = "BOOKED"
	SlotStatusBlocked   SlotStatus = "BLOCKED"
)

// SlotPeriod represents the day period categorisation used by legacy slot consumers.
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

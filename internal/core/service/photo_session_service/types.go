package photosessionservices

import (
	"time"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// UpdateSessionStatusInput contains data required to mutate a booking status.
type UpdateSessionStatusInput struct {
	SessionID      uint64
	PhotographerID uint64
	Status         string
}

// EnsureAgendaInput controls agenda bootstrap parameters.
type EnsureAgendaInput struct {
	PhotographerID uint64
	Timezone       string
}

// TimeOffInput represents the payload to create a time-off entry.
type TimeOffInput struct {
	PhotographerID uint64
	StartDate      time.Time
	EndDate        time.Time
	Reason         *string
	Location       *time.Location
}

// DeleteTimeOffInput represents the payload to remove a time-off entry.
type DeleteTimeOffInput struct {
	TimeOffID      uint64
	PhotographerID uint64
}

// UpdateTimeOffInput represents the payload to update a time-off entry.
type UpdateTimeOffInput struct {
	TimeOffID      uint64
	PhotographerID uint64
	StartDate      time.Time
	EndDate        time.Time
	Reason         *string
	Location       *time.Location
}

// ListTimeOffInput captures filters for time-off listing.
type ListTimeOffInput struct {
	PhotographerID uint64
	RangeFrom      time.Time
	RangeTo        time.Time
	Page           int
	Size           int
	Location       *time.Location
}

// TimeOffDetailInput carries identifiers to fetch a time-off entry.
type TimeOffDetailInput struct {
	TimeOffID      uint64
	PhotographerID uint64
	Timezone       string
}

// ListTimeOffOutput aggregates paginated time-off entries.
type ListTimeOffOutput struct {
	TimeOffs []photosessionmodel.AgendaEntryInterface
	Total    int64
	Page     int
	Size     int
	Timezone string
}

// TimeOffDetailResult represents a single time-off entry alongside timezone metadata.
type TimeOffDetailResult struct {
	TimeOff  photosessionmodel.AgendaEntryInterface
	Timezone string
}

// ListAgendaInput defines the input for listing agenda entries.
type ListAgendaInput struct {
	PhotographerID uint64
	StartDate      time.Time
	EndDate        time.Time
	Page           int
	Size           int
	Location       *time.Location
	EntryType      *photosessionmodel.AgendaEntryType
}

// ListAgendaOutput describes the agenda listing result.
type ListAgendaOutput struct {
	Slots    []AgendaSlot `json:"slots"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	Size     int          `json:"size"`
	Timezone string       `json:"timezone"`
}

// AgendaSlot represents an agenda entry rendered to clients.
type AgendaSlot struct {
	EntryID            uint64                              `json:"entryId"`
	PhotographerID     uint64                              `json:"photographerId"`
	EntryType          photosessionmodel.AgendaEntryType   `json:"entryType"`
	Source             photosessionmodel.AgendaEntrySource `json:"source"`
	SourceID           uint64                              `json:"sourceId,omitempty"`
	PhotoSessionID     *uint64                             `json:"photoSessionId,omitempty"`
	Start              time.Time                           `json:"start"`
	End                time.Time                           `json:"end"`
	Status             photosessionmodel.SlotStatus        `json:"status"`
	GroupID            string                              `json:"groupId"`
	IsHoliday          bool                                `json:"isHoliday"`
	IsTimeOff          bool                                `json:"isTimeOff"`
	HolidayLabels      []string                            `json:"holidayLabels,omitempty"`
	HolidayCalendarIDs []uint64                            `json:"holidayCalendarIds,omitempty"`
	Reason             string                              `json:"reason,omitempty"`
	Timezone           string                              `json:"timezone"`
}

// ListAvailabilityInput encapsulates range and pagination data for availability listing.
type ListAvailabilityInput struct {
	From      *time.Time
	To        *time.Time
	Page      int
	Size      int
	Sort      string
	Period    *photosessionmodel.SlotPeriod
	Location  *time.Location
	ListingID int64
}

// ListAvailabilityOutput aggregates computed availability slots.
type ListAvailabilityOutput struct {
	Slots    []AvailabilitySlot
	Total    int64
	Page     int
	Size     int
	Timezone string
}

// AvailabilitySlot represents a free window available for booking.
type AvailabilitySlot struct {
	SlotID         uint64
	PhotographerID uint64
	Start          time.Time
	End            time.Time
	Period         photosessionmodel.SlotPeriod
	SourceTimezone string
}

// ReserveSessionInput captures the necessary identifiers to reserve a photo session window.
type ReserveSessionInput struct {
	ListingID int64
	SlotID    uint64
	UserID    int64
}

// ReserveSessionOutput returns metadata about the reserved session.
type ReserveSessionOutput struct {
	PhotoSessionID uint64
	SlotID         uint64
	SlotStart      time.Time
	SlotEnd        time.Time
	PhotographerID uint64
	ListingID      int64
}

// ConfirmSessionInput holds data required to confirm a reserved session.
type ConfirmSessionInput struct {
	ListingID      int64
	PhotoSessionID uint64
	UserID         int64
}

// ConfirmSessionOutput reports the confirmed session metadata.
type ConfirmSessionOutput struct {
	PhotoSessionID uint64
	SlotStart      time.Time
	SlotEnd        time.Time
	PhotographerID uint64
	ListingID      int64
	Status         photosessionmodel.BookingStatus
}

// CancelSessionInput captures identifiers needed to cancel an existing session.
type CancelSessionInput struct {
	PhotoSessionID uint64
	UserID         int64
}

// CancelSessionOutput reports metadata about a cancelled session.
type CancelSessionOutput struct {
	PhotoSessionID uint64
	SlotStart      time.Time
	SlotEnd        time.Time
	PhotographerID uint64
	ListingID      int64
	ListingCode    uint32
}

// ListServiceAreasInput captures filters and pagination options when listing service areas.
type ListServiceAreasInput struct {
	PhotographerID uint64
	Page           int
	Size           int
}

// ListServiceAreasOutput bundles paginated service area entries.
type ListServiceAreasOutput struct {
	Areas []photosessionmodel.PhotographerServiceAreaInterface
	Total int64
	Page  int
	Size  int
}

// CreateServiceAreaInput represents the payload to create a new service area.
type CreateServiceAreaInput struct {
	PhotographerID uint64
	City           string
	State          string
}

// UpdateServiceAreaInput represents the payload to update an existing service area.
type UpdateServiceAreaInput struct {
	PhotographerID uint64
	ServiceAreaID  uint64
	City           string
	State          string
}

// DeleteServiceAreaInput carries identifiers to delete a service area.
type DeleteServiceAreaInput struct {
	PhotographerID uint64
	ServiceAreaID  uint64
}

// ServiceAreaDetailInput carries identifiers to retrieve a service area.
type ServiceAreaDetailInput struct {
	PhotographerID uint64
	ServiceAreaID  uint64
}

// ServiceAreaResult wraps a service area domain entity.
type ServiceAreaResult struct {
	Area photosessionmodel.PhotographerServiceAreaInterface
}

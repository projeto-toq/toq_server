package scheduleservices

import (
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
)

// CreateDefaultAgendaInput carries the data required to bootstrap a new agenda.
type CreateDefaultAgendaInput struct {
	ListingID int64
	OwnerID   int64
	Timezone  string
	ActorID   int64
}

// CreateBlockEntryInput captures the information to create a blocking entry.
type CreateBlockEntryInput struct {
	ListingID int64
	OwnerID   int64
	EntryType schedulemodel.EntryType
	StartsAt  time.Time
	EndsAt    time.Time
	Reason    string
	ActorID   int64
	Timezone  string
}

// UpdateBlockEntryInput represents the payload to update an existing blocking entry.
type UpdateBlockEntryInput struct {
	EntryID   uint64
	ListingID int64
	OwnerID   int64
	EntryType schedulemodel.EntryType
	StartsAt  time.Time
	EndsAt    time.Time
	Reason    string
	ActorID   int64
	Timezone  string
}

// DeleteEntryInput carries the identifiers required to remove an entry.
type DeleteEntryInput struct {
	EntryID   uint64
	ListingID int64
	OwnerID   int64
}

// AvailabilitySlot represents a continuous free window available for booking.
type AvailabilitySlot struct {
	StartsAt time.Time
	EndsAt   time.Time
}

// AvailabilityResult wraps paginated availability slots.
type AvailabilityResult struct {
	Slots    []AvailabilitySlot
	Total    int
	Timezone string
}

// FinishListingAgendaInput encapsulates data to finish the agenda creation workflow.
type FinishListingAgendaInput struct {
	ListingID int64
	OwnerID   int64
	ActorID   int64
}

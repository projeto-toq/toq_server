package listingservices

import (
	"time"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

const (
	defaultSlotsPage     = 1
	defaultSlotsPageSize = 20
	maxSlotsPageSize     = 100
	// reservationHoldTTL   = 15 * time.Minute
)

// ListPhotographerSlotsInput carries filtering and pagination data for slot listing.
type ListPhotographerSlotsInput struct {
	From              *time.Time
	To                *time.Time
	Period            *photosessionmodel.SlotPeriod
	Page              int
	Size              int
	Sort              string
	ListingIdentityID int64
	Location          *time.Location
}

// ListPhotographerSlotsOutput bundles slots and pagination metadata.
type ListPhotographerSlotsOutput struct {
	Slots []photosessionmodel.PhotographerSlotInterface
	Total int64
	Page  int
	Size  int
}

// ReservePhotoSessionInput holds identifiers needed to reserve a slot.
type ReservePhotoSessionInput struct {
	ListingIdentityID int64
	SlotID            uint64
}

// ReservePhotoSessionOutput returns metadata about the reserved slot.
type ReservePhotoSessionOutput struct {
	SlotID         uint64
	SlotStart      time.Time
	SlotEnd        time.Time
	PhotoSessionID uint64
	PhotographerID uint64
}

// CancelPhotoSessionInput identifies a scheduled photo session to cancel.
type CancelPhotoSessionInput struct {
	PhotoSessionID uint64
}

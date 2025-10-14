package listingservices

import (
	"time"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

const (
	defaultSlotsPage     = 1
	defaultSlotsPageSize = 20
	maxSlotsPageSize     = 100
	reservationHoldTTL   = 15 * time.Minute
)

// ListPhotographerSlotsInput carries filtering and pagination data for slot listing.
type ListPhotographerSlotsInput struct {
	From   *time.Time
	To     *time.Time
	Period *photosessionmodel.SlotPeriod
	Page   int
	Size   int
	Sort   string
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
	ListingID int64
	SlotID    uint64
}

// ReservePhotoSessionOutput returns reservation token and expiration timestamp.
type ReservePhotoSessionOutput struct {
	SlotID           uint64
	ReservationToken string
	ExpiresAt        time.Time
}

// ConfirmPhotoSessionInput encapsulates data to finalize a reservation into a booking.
type ConfirmPhotoSessionInput struct {
	ListingID        int64
	SlotID           uint64
	ReservationToken string
}

// ConfirmPhotoSessionOutput returns booking metadata after confirmation.
type ConfirmPhotoSessionOutput struct {
	PhotoSessionID uint64
	SlotID         uint64
	ScheduledStart time.Time
	ScheduledEnd   time.Time
}

// CancelPhotoSessionInput identifies a scheduled photo session to cancel.
type CancelPhotoSessionInput struct {
	PhotoSessionID uint64
}

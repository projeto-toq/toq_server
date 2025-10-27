package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// ToBookingModel converts a booking entity to its domain model representation.
func ToBookingModel(e entity.BookingEntity) photosessionmodel.PhotoSessionBookingInterface {
	model := photosessionmodel.NewPhotoSessionBooking()
	model.SetID(e.ID)
	model.SetSlotID(e.SlotID)
	model.SetListingID(e.ListingID)
	model.SetScheduledStart(e.ScheduledStart)
	model.SetScheduledEnd(e.ScheduledEnd)
	status, err := photosessionmodel.BookingStatusFromString(e.Status)
	if err == nil {
		model.SetStatus(status)
	} else {
		model.SetStatus(photosessionmodel.BookingStatusPendingApproval)
	}
	if e.Notes != nil {
		model.SetNotes(e.Notes)
	}
	return model
}

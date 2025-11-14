package converters

import (
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// ToBookingEntity maps a domain booking to its DB representation.
func ToBookingEntity(booking photosessionmodel.PhotoSessionBookingInterface) entity.Booking {
	reason := sql.NullString{}
	if val := booking.Reason(); val != nil {
		reason = sql.NullString{String: *val, Valid: true}
	}

	return entity.Booking{
		ID:                booking.ID(),
		AgendaEntryID:     booking.AgendaEntryID(),
		PhotographerID:    booking.PhotographerUserID(),
		ListingIdentityID: booking.ListingIdentityID(),
		StartsAt:          booking.StartsAt(),
		EndsAt:            booking.EndsAt(),
		Status:            string(booking.Status()),
		Reason:            reason,
	}
}

// ToBookingModel converts a DB entity into a domain booking model.
func ToBookingModel(entity entity.Booking) photosessionmodel.PhotoSessionBookingInterface {
	model := photosessionmodel.NewPhotoSessionBooking()
	model.SetID(entity.ID)
	model.SetAgendaEntryID(entity.AgendaEntryID)
	model.SetPhotographerUserID(entity.PhotographerID)
	model.SetListingIdentityID(entity.ListingIdentityID)
	model.SetStartsAt(entity.StartsAt)
	model.SetEndsAt(entity.EndsAt)
	model.SetStatus(photosessionmodel.BookingStatus(entity.Status))
	if entity.Reason.Valid {
		reason := entity.Reason.String
		model.SetReason(&reason)
	} else {
		model.SetReason(nil)
	}
	return model
}

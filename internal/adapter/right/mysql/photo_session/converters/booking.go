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

	createdAt := sql.NullTime{}
	if t, ok := booking.CreatedAt(); ok {
		createdAt = sql.NullTime{Time: t, Valid: true}
	}

	updatedAt := sql.NullTime{}
	if t, ok := booking.UpdatedAt(); ok {
		updatedAt = sql.NullTime{Time: t, Valid: true}
	}

	return entity.Booking{
		ID:             booking.ID(),
		AgendaEntryID:  booking.AgendaEntryID(),
		PhotographerID: booking.PhotographerUserID(),
		ListingID:      booking.ListingID(),
		StartsAt:       booking.StartsAt(),
		EndsAt:         booking.EndsAt(),
		Status:         string(booking.Status()),
		Reason:         reason,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

// ToBookingModel converts a DB entity into a domain booking model.
func ToBookingModel(entity entity.Booking) photosessionmodel.PhotoSessionBookingInterface {
	model := photosessionmodel.NewPhotoSessionBooking()
	model.SetID(entity.ID)
	model.SetAgendaEntryID(entity.AgendaEntryID)
	model.SetPhotographerUserID(entity.PhotographerID)
	model.SetListingID(entity.ListingID)
	model.SetStartsAt(entity.StartsAt)
	model.SetEndsAt(entity.EndsAt)
	model.SetStatus(photosessionmodel.BookingStatus(entity.Status))
	if entity.Reason.Valid {
		reason := entity.Reason.String
		model.SetReason(&reason)
	} else {
		model.SetReason(nil)
	}
	if entity.CreatedAt.Valid {
		model.SetCreatedAt(entity.CreatedAt.Time)
	}
	if entity.UpdatedAt.Valid {
		model.SetUpdatedAt(entity.UpdatedAt.Time)
	}
	return model
}

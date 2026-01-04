package converters

import (
	"database/sql"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// ToBookingEntity maps a domain booking to its DB representation, keeping nullable columns (reason, reservation_token)
// as sql.Null* and guarding against zeroed timestamps.
func ToBookingEntity(booking photosessionmodel.PhotoSessionBookingInterface) entity.Booking {
	reason := sql.NullString{}
	if val := booking.Reason(); val != nil {
		reason = sql.NullString{String: *val, Valid: true}
	}

	startsAt := sql.NullTime{}
	if ts := booking.StartsAt(); !ts.IsZero() {
		startsAt = sql.NullTime{Time: ts, Valid: true}
	}

	endsAt := sql.NullTime{}
	if ts := booking.EndsAt(); !ts.IsZero() {
		endsAt = sql.NullTime{Time: ts, Valid: true}
	}

	reservationToken := sql.NullString{}
	if token := booking.ReservationToken(); token != nil {
		reservationToken = sql.NullString{String: *token, Valid: true}
	}

	reservedUntil := sql.NullTime{}
	if ts := booking.ReservedUntil(); !ts.IsZero() {
		reservedUntil = sql.NullTime{Time: ts, Valid: true}
	}

	return entity.Booking{
		ID:                booking.ID(),
		AgendaEntryID:     booking.AgendaEntryID(),
		PhotographerID:    booking.PhotographerUserID(),
		ListingIdentityID: booking.ListingIdentityID(),
		StartsAt:          startsAt,
		EndsAt:            endsAt,
		Status:            string(booking.Status()),
		Reason:            reason,
		ReservationToken:  reservationToken,
		ReservedUntil:     reservedUntil,
	}
}

// ToBookingModel converts a DB entity into a domain booking model, applying zero values for NULL timestamps
// and nil for optional textual fields.
func ToBookingModel(entity entity.Booking) photosessionmodel.PhotoSessionBookingInterface {
	model := photosessionmodel.NewPhotoSessionBooking()
	model.SetID(entity.ID)
	model.SetAgendaEntryID(entity.AgendaEntryID)
	model.SetPhotographerUserID(entity.PhotographerID)
	model.SetListingIdentityID(entity.ListingIdentityID)

	if entity.StartsAt.Valid {
		model.SetStartsAt(entity.StartsAt.Time)
	} else {
		model.SetStartsAt(time.Time{})
	}

	if entity.EndsAt.Valid {
		model.SetEndsAt(entity.EndsAt.Time)
	} else {
		model.SetEndsAt(time.Time{})
	}

	model.SetStatus(photosessionmodel.BookingStatus(entity.Status))
	if entity.Reason.Valid {
		reason := entity.Reason.String
		model.SetReason(&reason)
	} else {
		model.SetReason(nil)
	}

	if entity.ReservationToken.Valid {
		token := entity.ReservationToken.String
		model.SetReservationToken(&token)
	} else {
		model.SetReservationToken(nil)
	}

	if entity.ReservedUntil.Valid {
		model.SetReservedUntil(entity.ReservedUntil.Time)
	} else {
		model.SetReservedUntil(time.Time{})
	}

	return model
}

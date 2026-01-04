package photosessionmodel

import "time"

type photoSessionBooking struct {
	id                uint64
	agendaEntryID     uint64
	photographerID    uint64
	listingIdentityID int64
	startsAt          time.Time
	endsAt            time.Time
	status            BookingStatus
	reason            *string
	reservationToken  *string
	reservedUntil     time.Time
	reservedValid     bool
}

func (b *photoSessionBooking) ID() uint64 { return b.id }

func (b *photoSessionBooking) SetID(id uint64) { b.id = id }

func (b *photoSessionBooking) AgendaEntryID() uint64 { return b.agendaEntryID }

func (b *photoSessionBooking) SetAgendaEntryID(id uint64) { b.agendaEntryID = id }

func (b *photoSessionBooking) PhotographerUserID() uint64 { return b.photographerID }

func (b *photoSessionBooking) SetPhotographerUserID(id uint64) { b.photographerID = id }

func (b *photoSessionBooking) ListingIdentityID() int64 { return b.listingIdentityID }

func (b *photoSessionBooking) SetListingIdentityID(id int64) { b.listingIdentityID = id }

func (b *photoSessionBooking) StartsAt() time.Time { return b.startsAt }

func (b *photoSessionBooking) SetStartsAt(value time.Time) { b.startsAt = value }

func (b *photoSessionBooking) EndsAt() time.Time { return b.endsAt }

func (b *photoSessionBooking) SetEndsAt(value time.Time) { b.endsAt = value }

func (b *photoSessionBooking) Status() BookingStatus { return b.status }

func (b *photoSessionBooking) SetStatus(status BookingStatus) { b.status = status }

func (b *photoSessionBooking) Reason() *string { return b.reason }

func (b *photoSessionBooking) SetReason(reason *string) { b.reason = reason }

func (b *photoSessionBooking) ReservationToken() *string { return b.reservationToken }

func (b *photoSessionBooking) SetReservationToken(token *string) { b.reservationToken = token }

func (b *photoSessionBooking) ReservedUntil() time.Time { return b.reservedUntil }

func (b *photoSessionBooking) SetReservedUntil(value time.Time) {
	b.reservedUntil = value
	b.reservedValid = true
}

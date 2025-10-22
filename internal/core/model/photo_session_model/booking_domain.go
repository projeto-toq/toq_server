package photosessionmodel

import "time"

type photoSessionBooking struct {
	id             uint64
	slotID         uint64
	listingID      int64
	scheduledStart time.Time
	scheduledEnd   time.Time
	status         BookingStatus
	notes          *string
}

func (b *photoSessionBooking) ID() uint64 { return b.id }

func (b *photoSessionBooking) SetID(id uint64) { b.id = id }

func (b *photoSessionBooking) SlotID() uint64 { return b.slotID }

func (b *photoSessionBooking) SetSlotID(id uint64) { b.slotID = id }

func (b *photoSessionBooking) ListingID() int64 { return b.listingID }

func (b *photoSessionBooking) SetListingID(id int64) { b.listingID = id }

func (b *photoSessionBooking) ScheduledStart() time.Time { return b.scheduledStart }

func (b *photoSessionBooking) SetScheduledStart(value time.Time) { b.scheduledStart = value }

func (b *photoSessionBooking) ScheduledEnd() time.Time { return b.scheduledEnd }

func (b *photoSessionBooking) SetScheduledEnd(value time.Time) { b.scheduledEnd = value }

func (b *photoSessionBooking) Status() BookingStatus { return b.status }

func (b *photoSessionBooking) SetStatus(status BookingStatus) { b.status = status }

func (b *photoSessionBooking) Notes() *string { return b.notes }

func (b *photoSessionBooking) SetNotes(notes *string) { b.notes = notes }

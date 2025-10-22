package photosessionmodel

import "time"

type photographerSlot struct {
	id                 uint64
	photographerUserID uint64
	slotDate           time.Time
	slotStart          time.Time
	slotEnd            time.Time
	period             SlotPeriod
	status             SlotStatus
	reservationToken   *string
	reservedUntil      *time.Time
	bookedAt           *time.Time
}

func (s *photographerSlot) ID() uint64 { return s.id }

func (s *photographerSlot) SetID(id uint64) { s.id = id }

func (s *photographerSlot) PhotographerUserID() uint64 { return s.photographerUserID }

func (s *photographerSlot) SetPhotographerUserID(id uint64) { s.photographerUserID = id }

func (s *photographerSlot) SlotDate() time.Time { return s.slotDate }

func (s *photographerSlot) SetSlotDate(date time.Time) { s.slotDate = date }

func (s *photographerSlot) SlotStart() time.Time { return s.slotStart }

func (s *photographerSlot) SetSlotStart(value time.Time) { s.slotStart = value }

func (s *photographerSlot) SlotEnd() time.Time { return s.slotEnd }

func (s *photographerSlot) SetSlotEnd(value time.Time) { s.slotEnd = value }

func (s *photographerSlot) Period() SlotPeriod { return s.period }

func (s *photographerSlot) SetPeriod(period SlotPeriod) { s.period = period }

func (s *photographerSlot) Status() SlotStatus { return s.status }

func (s *photographerSlot) SetStatus(status SlotStatus) { s.status = status }

func (s *photographerSlot) ReservationToken() *string { return s.reservationToken }

func (s *photographerSlot) SetReservationToken(token *string) { s.reservationToken = token }

func (s *photographerSlot) ReservedUntil() *time.Time { return s.reservedUntil }

func (s *photographerSlot) SetReservedUntil(value *time.Time) { s.reservedUntil = value }

func (s *photographerSlot) BookedAt() *time.Time { return s.bookedAt }

func (s *photographerSlot) SetBookedAt(value *time.Time) { s.bookedAt = value }

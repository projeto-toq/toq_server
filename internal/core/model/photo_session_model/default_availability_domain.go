package photosessionmodel

import "time"

type photographerDefaultAvailability struct {
	id                 uint64
	photographerUserID uint64
	weekday            time.Weekday
	period             SlotPeriod
	startHour          int
	slotsPerPeriod     int
	slotDurationMin    int
}

func (a *photographerDefaultAvailability) ID() uint64 { return a.id }

func (a *photographerDefaultAvailability) SetID(id uint64) { a.id = id }

func (a *photographerDefaultAvailability) PhotographerUserID() uint64 { return a.photographerUserID }

func (a *photographerDefaultAvailability) SetPhotographerUserID(id uint64) { a.photographerUserID = id }

func (a *photographerDefaultAvailability) Weekday() time.Weekday { return a.weekday }

func (a *photographerDefaultAvailability) SetWeekday(value time.Weekday) { a.weekday = value }

func (a *photographerDefaultAvailability) Period() SlotPeriod { return a.period }

func (a *photographerDefaultAvailability) SetPeriod(value SlotPeriod) { a.period = value }

func (a *photographerDefaultAvailability) StartHour() int { return a.startHour }

func (a *photographerDefaultAvailability) SetStartHour(value int) { a.startHour = value }

func (a *photographerDefaultAvailability) SlotsPerPeriod() int { return a.slotsPerPeriod }

func (a *photographerDefaultAvailability) SetSlotsPerPeriod(value int) { a.slotsPerPeriod = value }

func (a *photographerDefaultAvailability) SlotDurationMinutes() int { return a.slotDurationMin }

func (a *photographerDefaultAvailability) SetSlotDurationMinutes(value int) {
	a.slotDurationMin = value
}

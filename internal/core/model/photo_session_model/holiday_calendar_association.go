package photosessionmodel

import "time"

// HolidayCalendarAssociationInterface models the photographer_holiday_calendars table.
type HolidayCalendarAssociationInterface interface {
	ID() uint64
	SetID(id uint64)
	PhotographerUserID() uint64
	SetPhotographerUserID(id uint64)
	HolidayCalendarID() uint64
	SetHolidayCalendarID(id uint64)
	CreatedAt() (time.Time, bool)
	SetCreatedAt(t time.Time)
}

// NewHolidayCalendarAssociation returns a mutable association instance.
func NewHolidayCalendarAssociation() HolidayCalendarAssociationInterface {
	return &holidayCalendarAssociation{}
}

type holidayCalendarAssociation struct {
	id                 uint64
	photographerUserID uint64
	holidayCalendarID  uint64
	createdAt          time.Time
	createdAtDefined   bool
}

func (a *holidayCalendarAssociation) ID() uint64 {
	return a.id
}

func (a *holidayCalendarAssociation) SetID(id uint64) {
	a.id = id
}

func (a *holidayCalendarAssociation) PhotographerUserID() uint64 {
	return a.photographerUserID
}

func (a *holidayCalendarAssociation) SetPhotographerUserID(id uint64) {
	a.photographerUserID = id
}

func (a *holidayCalendarAssociation) HolidayCalendarID() uint64 {
	return a.holidayCalendarID
}

func (a *holidayCalendarAssociation) SetHolidayCalendarID(id uint64) {
	a.holidayCalendarID = id
}

func (a *holidayCalendarAssociation) CreatedAt() (time.Time, bool) {
	if !a.createdAtDefined {
		return time.Time{}, false
	}
	return a.createdAt, true
}

func (a *holidayCalendarAssociation) SetCreatedAt(t time.Time) {
	a.createdAt = t
	a.createdAtDefined = true
}

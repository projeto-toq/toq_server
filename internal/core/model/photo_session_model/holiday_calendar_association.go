package photosessionmodel

// HolidayCalendarAssociationInterface models the photographer_holiday_calendars table.
type HolidayCalendarAssociationInterface interface {
	ID() uint64
	SetID(id uint64)
	PhotographerUserID() uint64
	SetPhotographerUserID(id uint64)
	HolidayCalendarID() uint64
	SetHolidayCalendarID(id uint64)
}

// NewHolidayCalendarAssociation returns a mutable association instance.
func NewHolidayCalendarAssociation() HolidayCalendarAssociationInterface {
	return &holidayCalendarAssociation{}
}

type holidayCalendarAssociation struct {
	id                 uint64
	photographerUserID uint64
	holidayCalendarID  uint64
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

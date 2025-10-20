package holidaymodel

import "time"

type calendarDate struct {
	id          uint64
	calendarID  uint64
	holidayDate time.Time
	label       string
	recurrent   bool
}

func (d *calendarDate) ID() uint64 {
	return d.id
}

func (d *calendarDate) SetID(id uint64) {
	d.id = id
}

func (d *calendarDate) CalendarID() uint64 {
	return d.calendarID
}

func (d *calendarDate) SetCalendarID(id uint64) {
	d.calendarID = id
}

func (d *calendarDate) HolidayDate() time.Time {
	return d.holidayDate
}

func (d *calendarDate) SetHolidayDate(value time.Time) {
	d.holidayDate = value
}

func (d *calendarDate) Label() string {
	return d.label
}

func (d *calendarDate) SetLabel(value string) {
	d.label = value
}

func (d *calendarDate) IsRecurrent() bool {
	return d.recurrent
}

func (d *calendarDate) SetRecurrent(value bool) {
	d.recurrent = value
}

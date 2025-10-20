package holidaymodel

// CalendarInterface exposes the metadata stored for a holiday calendar.
type CalendarInterface interface {
	ID() uint64
	SetID(id uint64)
	Name() string
	SetName(value string)
	Scope() CalendarScope
	SetScope(value CalendarScope)
	State() (string, bool)
	SetState(value string)
	ClearState()
	CityIBGE() (string, bool)
	SetCityIBGE(value string)
	ClearCityIBGE()
	IsActive() bool
	SetActive(value bool)
	CreatedBy() int64
	SetCreatedBy(value int64)
	UpdatedBy() (int64, bool)
	SetUpdatedBy(value int64)
	ClearUpdatedBy()
}

// NewCalendar builds a new calendar entity.
func NewCalendar() CalendarInterface {
	return &calendar{}
}

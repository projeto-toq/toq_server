package converters

import (
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/holiday/entity"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
)

// ToCalendarModel converts CalendarEntity to domain interface.
func ToCalendarModel(e entity.CalendarEntity) holidaymodel.CalendarInterface {
	calendar := holidaymodel.NewCalendar()
	calendar.SetID(e.ID)
	calendar.SetName(e.Name)
	calendar.SetScope(holidaymodel.CalendarScope(e.Scope))
	calendar.SetTimezone(e.Timezone)
	if e.State.Valid {
		calendar.SetState(e.State.String)
	}
	if e.CityIBGE.Valid {
		calendar.SetCityIBGE(e.CityIBGE.String)
	}
	calendar.SetActive(e.IsActive)
	return calendar
}

// ToCalendarEntity converts domain calendar to persistence shape.
func ToCalendarEntity(model holidaymodel.CalendarInterface) entity.CalendarEntity {
	var state sql.NullString
	if value, ok := model.State(); ok {
		state = sql.NullString{String: value, Valid: true}
	}

	var city sql.NullString
	if value, ok := model.CityIBGE(); ok {
		city = sql.NullString{String: value, Valid: true}
	}

	return entity.CalendarEntity{
		ID:       model.ID(),
		Name:     model.Name(),
		Scope:    string(model.Scope()),
		State:    state,
		CityIBGE: city,
		IsActive: model.IsActive(),
		Timezone: model.Timezone(),
	}
}

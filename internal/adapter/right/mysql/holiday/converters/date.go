package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/holiday/entity"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
)

// ToDateModel converts DateEntity to domain interface.
func ToDateModel(e entity.DateEntity) holidaymodel.CalendarDateInterface {
	date := holidaymodel.NewCalendarDate()
	date.SetID(e.ID)
	date.SetCalendarID(e.CalendarID)
	date.SetHolidayDate(e.Holiday)
	date.SetLabel(e.Label)
	date.SetRecurrent(e.Recurrent)
	return date
}

// ToDateEntity converts domain date to persistence shape.
func ToDateEntity(model holidaymodel.CalendarDateInterface) entity.DateEntity {
	return entity.DateEntity{
		ID:         model.ID(),
		CalendarID: model.CalendarID(),
		Holiday:    model.HolidayDate(),
		Label:      model.Label(),
		Recurrent:  model.IsRecurrent(),
	}
}

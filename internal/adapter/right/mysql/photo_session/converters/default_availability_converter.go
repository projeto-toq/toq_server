package converters

import (
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// ToDefaultAvailabilityModel converts entity to domain model.
func ToDefaultAvailabilityModel(e entity.DefaultAvailabilityEntity) photosessionmodel.PhotographerDefaultAvailabilityInterface {
	record := photosessionmodel.NewPhotographerDefaultAvailability()
	record.SetID(e.ID)
	record.SetPhotographerUserID(e.PhotographerUserID)
	record.SetWeekday(time.Weekday(e.Weekday))
	record.SetPeriod(photosessionmodel.SlotPeriod(e.Period))
	record.SetStartHour(e.StartHour)
	record.SetSlotsPerPeriod(e.SlotsPerPeriod)
	record.SetSlotDurationMinutes(e.SlotDurationMin)
	return record
}

// FromDefaultAvailabilityModel converts domain model to entity representation.
func FromDefaultAvailabilityModel(model photosessionmodel.PhotographerDefaultAvailabilityInterface) entity.DefaultAvailabilityEntity {
	return entity.DefaultAvailabilityEntity{
		ID:                 model.ID(),
		PhotographerUserID: model.PhotographerUserID(),
		Weekday:            int(model.Weekday()),
		Period:             string(model.Period()),
		StartHour:          model.StartHour(),
		SlotsPerPeriod:     model.SlotsPerPeriod(),
		SlotDurationMin:    model.SlotDurationMinutes(),
	}
}

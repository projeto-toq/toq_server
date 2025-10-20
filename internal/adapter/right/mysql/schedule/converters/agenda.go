package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entity"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
)

// ToAgendaModel converts an AgendaEntity to the domain interface.
func ToAgendaModel(e entity.AgendaEntity) schedulemodel.AgendaInterface {
	agenda := schedulemodel.NewAgenda()
	agenda.SetID(e.ID)
	agenda.SetListingID(e.ListingID)
	agenda.SetOwnerID(e.OwnerID)
	agenda.SetTimezone(e.Timezone)
	return agenda
}

// ToAgendaEntity converts the domain object into its persistence shape.
func ToAgendaEntity(model schedulemodel.AgendaInterface) entity.AgendaEntity {
	return entity.AgendaEntity{
		ID:        model.ID(),
		ListingID: model.ListingID(),
		OwnerID:   model.OwnerID(),
		Timezone:  model.Timezone(),
	}
}

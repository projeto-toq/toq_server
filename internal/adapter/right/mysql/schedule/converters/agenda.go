// Package scheduleconverters translates between MySQL entities and schedule domain models, keeping persistence concerns isolated.
// Functions here must be pure data mappers without side effects or business rules (guide Section 8).
package scheduleconverters

import (
	scheduleentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entities"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
)

// AgendaEntityToDomain converts AgendaEntity into the domain representation.
func AgendaEntityToDomain(e scheduleentity.AgendaEntity) schedulemodel.AgendaInterface {
	agenda := schedulemodel.NewAgenda()
	agenda.SetID(e.ID)
	agenda.SetListingIdentityID(e.ListingIdentityID)
	agenda.SetOwnerID(e.OwnerID)
	agenda.SetTimezone(e.Timezone)
	return agenda
}

// AgendaDomainToEntity converts a domain agenda into its persistence shape.
func AgendaDomainToEntity(model schedulemodel.AgendaInterface) scheduleentity.AgendaEntity {
	return scheduleentity.AgendaEntity{
		ID:                model.ID(),
		ListingIdentityID: model.ListingIdentityID(),
		OwnerID:           model.OwnerID(),
		Timezone:          model.Timezone(),
	}
}

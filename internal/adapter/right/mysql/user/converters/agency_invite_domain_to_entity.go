package userconverters

import (
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// AgencyInviteDomainToEntity converts a domain model to a database entity
//
// This converter handles the translation from clean domain types to database-ready
// entity for insertion/update operations on the agency_invites table.
//
// Conversion Rules:
//   - All fields are NOT NULL per schema (no sql.Null* needed)
//   - ID may be 0 for new records (populated by AUTO_INCREMENT)
//   - PhoneNumber must be in E.164 format (validated by service layer)
//
// Parameters:
//   - domain: InviteInterface from core layer
//
// Returns:
//   - entity: AgencyInvite ready for database operations
//
// Example:
//
//	invite := usermodel.NewInvite()
//	invite.SetAgencyID(123)
//	invite.SetPhoneNumber("+5511999999999")
//	entity := AgencyInviteDomainToEntity(invite)
//	// entity ready for INSERT query
func AgencyInviteDomainToEntity(domain usermodel.InviteInterface) (entity userentity.AgencyInvite) {
	entity = userentity.AgencyInvite{}
	entity.ID = uint32(domain.GetID())
	entity.AgencyID = uint32(domain.GetAgencyID())
	entity.PhoneNumber = domain.GetPhoneNumber()

	return
}

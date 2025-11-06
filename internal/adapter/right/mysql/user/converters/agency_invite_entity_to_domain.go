package userconverters

import (
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// AgencyInviteEntityToDomainTyped converts a type-safe AgencyInvite to InviteInterface domain model.
//
// This is the preferred converter function as it uses compile-time type safety instead of runtime
// type assertions. All fields in agency_invites are NOT NULL per schema, so no null handling needed.
//
// Parameters:
//   - entity: Strongly-typed AgencyInvite with all non-nullable fields
//
// Returns:
//   - InviteInterface: Domain model with appropriate getters/setters
//
// Example:
//
//	entity := userentity.AgencyInvite{
//	    ID: 123,
//	    AgencyID: 456,
//	    PhoneNumber: "+5511999999999",
//	}
//	invite := AgencyInviteEntityToDomainTyped(entity)
func AgencyInviteEntityToDomainTyped(entity userentity.AgencyInvite) usermodel.InviteInterface {
	domain := usermodel.NewInvite()

	domain.SetID(entity.ID)
	domain.SetAgencyID(entity.AgencyID)
	domain.SetPhoneNumber(entity.PhoneNumber)

	return domain
}

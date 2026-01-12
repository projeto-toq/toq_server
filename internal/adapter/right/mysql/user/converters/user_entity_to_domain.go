package userconverters

import (
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// UserEntityToDomain converts a database entity to a domain model
//
// This converter handles the translation from database-specific types (sql.Null*)
// to clean domain types, ensuring the core layer remains decoupled from database concerns.
//
// Conversion Rules:
//   - sql.NullString → string (empty string if NULL)
//   - sql.NullTime → time.Time (zero time if NULL)
//   - TINYINT(1) → bool (true if non-zero)
//
// Parameters:
//   - entity: UserEntity from database query
//
// Returns:
//   - user: UserInterface with all fields populated from entity
//
// Note: ActiveRole is NOT set here - must be populated by Permission Service
//
//	This maintains separation of concerns between User and Permission domains
func UserEntityToDomain(entity userentity.UserEntity) usermodel.UserInterface {
	user := usermodel.NewUser()

	// Map mandatory fields (NOT NULL in schema)
	user.SetID(int64(entity.ID))
	user.SetFullName(entity.FullName)
	user.SetNationalID(entity.NationalID)
	user.SetBornAt(entity.BornAt)
	user.SetPhoneNumber(entity.PhoneNumber)
	user.SetEmail(entity.Email)
	user.SetZipCode(entity.ZipCode)
	user.SetStreet(entity.Street)
	user.SetNumber(entity.Number)
	user.SetNeighborhood(entity.Neighborhood)
	user.SetCity(entity.City)
	user.SetState(entity.State)
	user.SetPassword(entity.Password)
	user.SetOptStatus(entity.OptStatus)
	user.SetLastActivityAt(entity.LastActivityAt)
	user.SetCreatedAt(entity.CreatedAt)
	user.SetDeleted(entity.Deleted)

	// Map optional fields (NULL in schema) - check Valid before accessing
	if entity.NickName.Valid {
		user.SetNickName(entity.NickName.String)
	}

	if entity.CreciNumber.Valid {
		user.SetCreciNumber(entity.CreciNumber.String)
	}

	if entity.CreciState.Valid {
		user.SetCreciState(entity.CreciState.String)
	}

	if entity.CreciValidity.Valid {
		user.SetCreciValidity(entity.CreciValidity.Time)
	}

	if entity.Complement.Valid {
		user.SetComplement(entity.Complement.String)
	}

	// NEW: Map blocking fields
	if entity.BlockedUntil.Valid {
		t := entity.BlockedUntil.Time
		user.SetBlockedUntil(&t)
	} else {
		user.SetBlockedUntil(nil)
	}

	user.SetPermanentlyBlocked(entity.PermanentlyBlocked)

	return user
}

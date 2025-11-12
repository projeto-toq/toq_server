package userconverters

import (
	"database/sql"

	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// UserDomainToEntity converts a domain model to a database entity
//
// This converter handles the translation from clean domain types to database-specific
// types (sql.Null*), preparing data for database insertion/update.
//
// Conversion Rules:
//   - string → sql.NullString (Valid=true if non-empty)
//   - time.Time → sql.NullTime (Valid=true if not zero time)
//   - bool → TINYINT(1) (true = 1, false = 0)
//
// Parameters:
//   - domain: UserInterface from core layer
//
// Returns:
//   - entity: UserEntity ready for database operations
//
// Important:
//   - ID may be 0 for new records (populated by AUTO_INCREMENT)
//   - Empty strings are converted to NULL for optional fields
//   - Zero times (IsZero()) are converted to NULL for optional date fields
//   - Photo field is not set here (managed separately by photo upload flow)
func UserDomainToEntity(domain usermodel.UserInterface) (entity userentity.UserEntity) {
	entity = userentity.UserEntity{}

	// Map mandatory fields (NOT NULL in schema)
	entity.ID = uint32(domain.GetID())
	entity.FullName = domain.GetFullName()
	entity.NationalID = domain.GetNationalID()
	entity.BornAt = domain.GetBornAt()
	entity.PhoneNumber = domain.GetPhoneNumber()
	entity.Email = domain.GetEmail()
	entity.ZipCode = domain.GetZipCode()
	entity.Street = domain.GetStreet()
	entity.Number = domain.GetNumber()
	entity.Neighborhood = domain.GetNeighborhood()
	entity.City = domain.GetCity()
	entity.State = domain.GetState()
	entity.Password = domain.GetPassword()
	entity.OptStatus = domain.IsOptStatus()
	entity.LastActivityAt = domain.GetLastActivityAt()
	entity.Deleted = domain.IsDeleted()

	// Map optional fields - convert to sql.Null* with Valid based on value presence
	nickName := domain.GetNickName()
	entity.NickName = sql.NullString{
		String: nickName,
		Valid:  nickName != "",
	}

	creciNumber := domain.GetCreciNumber()
	entity.CreciNumber = sql.NullString{
		String: creciNumber,
		Valid:  creciNumber != "",
	}

	creciState := domain.GetCreciState()
	entity.CreciState = sql.NullString{
		String: creciState,
		Valid:  creciState != "",
	}

	creciValidity := domain.GetCreciValidity()
	entity.CreciValidity = sql.NullTime{
		Time:  creciValidity,
		Valid: !creciValidity.IsZero(),
	}

	complement := domain.GetComplement()
	entity.Complement = sql.NullString{
		String: complement,
		Valid:  complement != "",
	}

	// NEW: Map blocking fields
	blockedUntil := domain.GetBlockedUntil()
	if blockedUntil != nil {
		entity.BlockedUntil = sql.NullTime{
			Time:  *blockedUntil,
			Valid: true,
		}
	} else {
		entity.BlockedUntil = sql.NullTime{
			Valid: false,
		}
	}

	entity.PermanentlyBlocked = domain.IsPermanentlyBlocked()

	return
}

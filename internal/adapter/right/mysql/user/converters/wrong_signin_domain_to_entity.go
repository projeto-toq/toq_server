package userconverters

import (
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// WrongSignInDomainToEntity converts a domain model to a database entity
//
// This converter handles the translation from domain types to database entity
// for the temp_wrong_signin table used in brute-force prevention.
//
// Conversion Rules:
//   - int64 → uint8 for FailedAttempts (valid range 0-255)
//   - time.Time remains unchanged (NOT NULL field)
//   - All fields are required (no NULL handling needed)
//
// Parameters:
//   - domain: WrongSigninInterface from core layer
//
// Returns:
//   - entity: WrongSignInEntity ready for database operations
//
// Important:
//   - UserID must reference existing user (foreign key constraint)
//   - FailedAttempts max value is 255 (TINYINT UNSIGNED limit)
//   - LastAttemptAt should be set to current time on each failure
//
// Example:
//
//	wrongSignin := usermodel.NewWrongSignin()
//	wrongSignin.SetUserID(123)
//	wrongSignin.SetFailedAttempts(3)
//	wrongSignin.SetLastAttemptAt(time.Now())
//	entity := WrongSignInDomainToEntity(wrongSignin)
func WrongSignInDomainToEntity(domain usermodel.WrongSigninInterface) (entity userentity.WrongSignInEntity) {
	entity = userentity.WrongSignInEntity{}
	entity.UserID = uint32(domain.GetUserID())
	entity.FailedAttempts = uint8(domain.GetFailedAttempts())
	entity.LastAttemptAT = domain.GetLastAttemptAt()
	return
}

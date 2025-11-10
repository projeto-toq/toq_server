package userconverters

import (
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// WrongSignInEntityToDomainTyped converts a type-safe WrongSignInEntity to WrongSigninInterface domain model.
//
// This is the preferred converter function as it uses compile-time type safety instead of runtime
// type assertions. The LastAttemptAT field (note: AT not At) is always non-null per schema.
//
// Parameters:
//   - entity: Strongly-typed WrongSignInEntity with time.Time for timestamp
//
// Returns:
//   - WrongSigninInterface: Domain model with appropriate getters/setters
//
// Example:
//
//	entity := userentity.WrongSignInEntity{
//	    UserID: 123,
//	    FailedAttempts: 3,
//	    LastAttemptAT: time.Now(),
//	}
//	wrongSignin := WrongSignInEntityToDomainTyped(entity)
func WrongSignInEntityToDomainTyped(entity userentity.WrongSignInEntity) usermodel.WrongSigninInterface {
	wsi := usermodel.NewWrongSignin()

	wsi.SetUserID(int64(entity.UserID))
	wsi.SetFailedAttempts(int64(entity.FailedAttempts)) // uint8 -> int64 conversion
	wsi.SetLastAttemptAt(entity.LastAttemptAT)

	return wsi
}

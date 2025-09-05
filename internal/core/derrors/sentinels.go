package derrors

import "errors"

// Domain sentinel errors (use errors.Is to check)
var (
	ErrEmailChangeNotPending  = errors.New("email change not pending")
	ErrEmailChangeCodeInvalid = errors.New("invalid email change code")
	ErrEmailChangeCodeExpired = errors.New("email change code expired")
	ErrEmailAlreadyInUse      = errors.New("email already in use")
	ErrSameEmailAsCurrent     = errors.New("new email is the same as current")

	ErrPhoneChangeNotPending  = errors.New("phone change not pending")
	ErrPhoneChangeCodeInvalid = errors.New("invalid phone change code")
	ErrPhoneChangeCodeExpired = errors.New("phone change code expired")
	ErrPhoneAlreadyInUse      = errors.New("phone already in use")
	ErrSamePhoneAsCurrent     = errors.New("new phone is the same as current")

	ErrUserActiveRoleMissing = errors.New("active role missing for user")
)

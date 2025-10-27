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

	ErrCPFInvalid          = errors.New("cpf invalid")
	ErrCPFBirthDateInvalid = errors.New("cpf birth date invalid")
	ErrCPFNotFound         = errors.New("cpf not found")
	ErrCPFStatusIrregular  = errors.New("cpf status irregular")
	ErrCPFDataMismatch     = errors.New("cpf data mismatch")

	ErrCNPJInvalid  = errors.New("cnpj invalid")
	ErrCNPJNotFound = errors.New("cnpj not found")

	ErrCEPInvalid  = errors.New("invalid CEP")
	ErrCEPNotFound = errors.New("CEP not found")

	ErrSlotUnavailable           = errors.New("photographer slot unavailable")
	ErrReservationExpired        = errors.New("photographer slot reservation expired")
	ErrListingNotEligible        = errors.New("listing not eligible for photo session")
	ErrPhotoSessionNotCancelable = errors.New("photo session cannot be cancelled")
	ErrPhotoSessionPending       = errors.New("photo session awaiting photographer decision")
	ErrPhotoSessionAlreadyFinal  = errors.New("photo session already finalized")

	ErrRoleNotSystem          = errors.New("role is not marked as system role")
	ErrAdminRoleProtected     = errors.New("admin role cannot be altered")
	ErrCannotDeleteLoggedUser = errors.New("logged user cannot delete itself")
	ErrUserAlreadyDeleted     = errors.New("user already deleted")
	ErrRoleSlugImmutable      = errors.New("role slug cannot be changed")
	ErrSystemUserRoleMismatch = errors.New("user role slug mismatch")
	ErrRoleDeletionHasUsers   = errors.New("role still assigned to users")
)

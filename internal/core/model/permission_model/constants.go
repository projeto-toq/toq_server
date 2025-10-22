package permissionmodel

import "regexp"

// RoleSlug representa os slugs de roles do sistema
type RoleSlug string

const (
	RoleSlugRoot             RoleSlug = "root"
	RoleSlugOwner            RoleSlug = "owner"
	RoleSlugRealtor          RoleSlug = "realtor"
	RoleSlugAgency           RoleSlug = "agency"
	RoleSlugPhotographer     RoleSlug = "photographer"
	RoleSlugAttendantRealtor RoleSlug = "attendantRealtor"
	RoleSlugAttendantOwner   RoleSlug = "attendantOwner"
	RoleSlugAttendant        RoleSlug = "attendant"
	RoleSlugManager          RoleSlug = "manager"
)

// String implementa fmt.Stringer
func (rs RoleSlug) String() string {
	return string(rs)
}

// IsValid verifica se o slug é válido
var roleSlugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{1,63}$`)

func (rs RoleSlug) IsValid() bool {
	if rs == "" {
		return false
	}
	return roleSlugPattern.MatchString(rs.String())
}

// GetAllRoleSlugs retorna todos os slugs válidos
func GetAllRoleSlugs() []RoleSlug {
	return []RoleSlug{
		RoleSlugRoot,
		RoleSlugOwner,
		RoleSlugRealtor,
		RoleSlugAgency,
		RoleSlugPhotographer,
		RoleSlugAttendantRealtor,
		RoleSlugAttendantOwner,
		RoleSlugAttendant,
		RoleSlugManager,
	}
}

// UserRoleStatus representa os possíveis status de um user_role
type UserRoleStatus int

const (
	StatusActive          UserRoleStatus = iota // normal user status 0
	StatusBlocked                               // blocked by admin 1
	StatusTempBlocked                           // temporarily blocked due to failed signin attempts 2
	StatusPendingBoth                           // awaiting both email and phone confirmation 3
	StatusPendingEmail                          // awaiting email confirmation 4
	StatusPendingPhone                          // awaiting phone confirmation 5
	StatusPendingCreci                          // awaiting creci images to be uploaded 6
	StatusPendingCnpj                           // awaiting cnpj images to be uploaded 7
	StatusPendingManual                         // awaiting manual verification by admin 8
	StatusRejected                              // admin reject the documentation (legacy/general) 9
	StatusRefusedImage                          // refused due to image issues (e.g., unreadable/invalid) 10
	StatusRefusedDocument                       // refused due to document mismatch/invalidity 11
	StatusRefusedData                           // refused due to data inconsistency 12
	StatusDeleted                               // user request the deletion of the account 13
)

// String implementa fmt.Stringer para UserRoleStatus
func (us UserRoleStatus) String() string {
	statuses := [...]string{
		"active",
		"blocked",
		"temp_blocked",
		"pending_both",
		"pending_email",
		"pending_phone",
		"pending_creci",
		"pending_cnpj",
		"pending_manual",
		"rejected",
		"refused_image",
		"refused_document",
		"refused_data",
		"deleted",
	}
	if us < StatusActive || int(us) >= len(statuses) {
		return "unknown"
	}
	return statuses[us]
}

// IsManualApprovalTarget verifies if the status is allowed for manual approval actions.
func IsManualApprovalTarget(status UserRoleStatus) bool {
	switch status {
	case StatusActive, StatusRejected, StatusRefusedImage, StatusRefusedDocument, StatusRefusedData:
		return true
	default:
		return false
	}
}

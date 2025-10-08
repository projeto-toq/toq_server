package permissionmodel

// RoleSlug representa os slugs de roles do sistema
type RoleSlug string

const (
	RoleSlugRoot    RoleSlug = "root"
	RoleSlugOwner   RoleSlug = "owner"
	RoleSlugRealtor RoleSlug = "realtor"
	RoleSlugAgency  RoleSlug = "agency"
)

// String implementa fmt.Stringer
func (rs RoleSlug) String() string {
	return string(rs)
}

// IsValid verifica se o slug é válido
func (rs RoleSlug) IsValid() bool {
	switch rs {
	case RoleSlugRoot, RoleSlugOwner, RoleSlugRealtor, RoleSlugAgency:
		return true
	default:
		return false
	}
}

// GetAllRoleSlugs retorna todos os slugs válidos
func GetAllRoleSlugs() []RoleSlug {
	return []RoleSlug{
		RoleSlugRoot,
		RoleSlugOwner,
		RoleSlugRealtor,
		RoleSlugAgency,
	}
}

// PermissionResource representa os recursos do sistema
type PermissionResource string

const (
	ResourceUser       PermissionResource = "user"
	ResourceListing    PermissionResource = "listing"
	ResourceVisit      PermissionResource = "visit"
	ResourceOffer      PermissionResource = "offer"
	ResourceRole       PermissionResource = "role"
	ResourcePermission PermissionResource = "permission"
)

// PermissionAction representa as ações disponíveis
type PermissionAction string

const (
	ActionCreate PermissionAction = "create"
	ActionRead   PermissionAction = "read"
	ActionUpdate PermissionAction = "update"
	ActionDelete PermissionAction = "delete"
	ActionAssign PermissionAction = "assign"
	ActionRevoke PermissionAction = "revoke"
)

// String implementa fmt.Stringer
func (pr PermissionResource) String() string {
	return string(pr)
}

func (pa PermissionAction) String() string {
	return string(pa)
}

// IsValidResource verifica se o resource é válido
func (pr PermissionResource) IsValid() bool {
	switch pr {
	case ResourceUser, ResourceListing, ResourceVisit, ResourceOffer, ResourceRole, ResourcePermission:
		return true
	default:
		return false
	}
}

// IsValidAction verifica se a action é válida
func (pa PermissionAction) IsValid() bool {
	switch pa {
	case ActionCreate, ActionRead, ActionUpdate, ActionDelete, ActionAssign, ActionRevoke:
		return true
	default:
		return false
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
	StatusInvitePending                         // realtor was invited and is pending acceptance 14
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
		"invite_pending",
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

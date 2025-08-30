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

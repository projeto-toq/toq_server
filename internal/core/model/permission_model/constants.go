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

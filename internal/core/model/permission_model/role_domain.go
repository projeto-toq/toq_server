package permissionmodel

// role representa um papel/função no sistema
type role struct {
	id           int64
	name         string
	slug         string
	description  string
	isSystemRole bool
	isActive     bool
}

func NewRole() RoleInterface {
	return &role{
		isActive: true,
	}
}

func (r *role) GetID() int64 {
	return r.id
}

func (r *role) SetID(id int64) {
	r.id = id
}

func (r *role) GetName() string {
	return r.name
}

func (r *role) SetName(name string) {
	r.name = name
}

func (r *role) GetSlug() string {
	return r.slug
}

func (r *role) SetSlug(slug string) {
	r.slug = slug
}

func (r *role) GetDescription() string {
	return r.description
}

func (r *role) SetDescription(description string) {
	r.description = description
}

func (r *role) GetIsSystemRole() bool {
	return r.isSystemRole
}

func (r *role) SetIsSystemRole(isSystemRole bool) {
	r.isSystemRole = isSystemRole
}

func (r *role) GetIsActive() bool {
	return r.isActive
}

func (r *role) SetIsActive(isActive bool) {
	r.isActive = isActive
}

package permissionmodel

type RoleInterface interface {
	GetID() int64
	SetID(id int64)
	GetName() string
	SetName(name string)
	GetSlug() string
	SetSlug(slug string)
	GetDescription() string
	SetDescription(description string)
	GetIsSystemRole() bool
	SetIsSystemRole(isSystemRole bool)
	GetIsActive() bool
	SetIsActive(isActive bool)
}

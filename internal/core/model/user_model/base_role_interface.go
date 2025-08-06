package usermodel

type BaseRoleInterface interface {
	GetID() int64
	SetID(id int64)
	GetRole() UserRole
	SetRole(role UserRole)
	GetName() string
	SetName(name string)
	GetPrivileges() []PrivilegeInterface
	SetPrivileges(privileges []PrivilegeInterface)
	AddPrivilege(privilege PrivilegeInterface)
	RemovePrivilege(privilege PrivilegeInterface)
}

func NewBaseRole() BaseRoleInterface {
	return &baseRole{}
}

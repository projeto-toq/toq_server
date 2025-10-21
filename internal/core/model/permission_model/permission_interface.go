package permissionmodel

type PermissionInterface interface {
	GetID() int64
	SetID(id int64)
	GetName() string
	SetName(name string)
	GetDescription() string
	SetDescription(description string)
	GetAction() string
	SetAction(action string)
	GetIsActive() bool
	SetIsActive(isActive bool)
}

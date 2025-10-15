package permissionmodel

type PermissionInterface interface {
	GetID() int64
	SetID(id int64)
	GetName() string
	SetName(name string)
	GetDescription() string
	SetDescription(description string)
	GetResource() string
	SetResource(resource string)
	GetAction() string
	SetAction(action string)
	GetConditions() map[string]interface{}
	SetConditions(conditions map[string]interface{})
	SetConditionsFromJSON(jsonData []byte) error
	GetConditionsAsJSON() ([]byte, error)
	GetIsActive() bool
	SetIsActive(isActive bool)
}

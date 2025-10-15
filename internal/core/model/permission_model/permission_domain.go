package permissionmodel

import (
	"encoding/json"
)

type permission struct {
	id          int64
	name        string
	description string
	resource    string
	action      string
	conditions  map[string]interface{}
	isActive    bool
}

func NewPermission() PermissionInterface {
	return &permission{
		isActive: true,
	}
}

func (p *permission) GetID() int64 {
	return p.id
}

func (p *permission) SetID(id int64) {
	p.id = id
}

func (p *permission) GetName() string {
	return p.name
}

func (p *permission) SetName(name string) {
	p.name = name
}

func (p *permission) GetDescription() string {
	return p.description
}

func (p *permission) SetDescription(description string) {
	p.description = description
}

func (p *permission) GetResource() string {
	return p.resource
}

func (p *permission) SetResource(resource string) {
	p.resource = resource
}

func (p *permission) GetAction() string {
	return p.action
}

func (p *permission) SetAction(action string) {
	p.action = action
}

func (p *permission) GetConditions() map[string]interface{} {
	return p.conditions
}

func (p *permission) SetConditions(conditions map[string]interface{}) {
	p.conditions = conditions
}

func (p *permission) GetIsActive() bool {
	return p.isActive
}

func (p *permission) SetIsActive(isActive bool) {
	p.isActive = isActive
}

func (p *permission) SetConditionsFromJSON(jsonData []byte) error {
	if len(jsonData) == 0 {
		p.conditions = nil
		return nil
	}

	var conditions map[string]interface{}
	err := json.Unmarshal(jsonData, &conditions)
	if err != nil {
		return err
	}

	p.conditions = conditions
	return nil
}

func (p *permission) GetConditionsAsJSON() ([]byte, error) {
	if p.conditions == nil {
		return nil, nil
	}

	return json.Marshal(p.conditions)
}

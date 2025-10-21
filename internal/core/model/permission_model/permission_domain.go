package permissionmodel

// permission represents a simplified permission focused on HTTP actions.
type permission struct {
	id          int64
	name        string
	description string
	action      string
	isActive    bool
}

// NewPermission creates a new permission with default active flag.
func NewPermission() PermissionInterface {
	return &permission{isActive: true}
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

func (p *permission) GetAction() string {
	return p.action
}

func (p *permission) SetAction(action string) {
	p.action = action
}

func (p *permission) GetIsActive() bool {
	return p.isActive
}

func (p *permission) SetIsActive(isActive bool) {
	p.isActive = isActive
}

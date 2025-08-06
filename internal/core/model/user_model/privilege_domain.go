package usermodel

type privilege struct {
	id      int64
	role_id int64
	service GRPCService
	method  uint8
	allowed bool
}

func (p *privilege) ID() int64 {
	return p.id
}

func (p *privilege) SetID(id int64) {
	p.id = id
}

func (p *privilege) RoleID() int64 {
	return p.role_id
}

func (p *privilege) SetRoleID(role_id int64) {
	p.role_id = role_id
}

func (p *privilege) Service() GRPCService {
	return p.service
}

func (p *privilege) SetService(service GRPCService) {
	p.service = service
}

func (p *privilege) Method() uint8 {
	return p.method
}

func (p *privilege) SetMethod(method uint8) {
	p.method = method
}

func (p *privilege) Allowed() bool {
	return p.allowed
}

func (p *privilege) SetAllowed(allowed bool) {
	p.allowed = allowed
}

package usermodel

type PrivilegeInterface interface {
	ID() int64
	SetID(id int64)
	RoleID() int64
	SetRoleID(role_id int64)
	Service() GRPCService
	SetService(service GRPCService)
	Method() uint8
	SetMethod(rpc uint8)
	Allowed() bool
	SetAllowed(isAllowed bool)
}

func NewPrivilege() PrivilegeInterface {
	return &privilege{}
}

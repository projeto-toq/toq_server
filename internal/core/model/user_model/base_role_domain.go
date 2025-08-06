package usermodel

type baseRole struct {
	id         int64
	role       UserRole
	name       string
	privileges []PrivilegeInterface
}

func (br *baseRole) GetID() int64 {
	return br.id
}

func (br *baseRole) SetID(id int64) {
	br.id = id
}

func (br *baseRole) GetRole() UserRole {
	return br.role
}

func (br *baseRole) SetRole(role UserRole) {
	br.role = role
}

func (br *baseRole) GetName() string {
	return br.name
}

func (br *baseRole) SetName(name string) {
	br.name = name
}

func (br *baseRole) GetPrivileges() []PrivilegeInterface {
	return br.privileges
}

func (br *baseRole) SetPrivileges(privileges []PrivilegeInterface) {
	br.privileges = privileges
}

func (br *baseRole) AddPrivilege(privilege PrivilegeInterface) {
	br.privileges = append(br.privileges, privilege)
}

func (br *baseRole) RemovePrivilege(privilege PrivilegeInterface) {
	for i, p := range br.privileges {
		if p.ID() == privilege.ID() {
			br.privileges = append(br.privileges[:i], br.privileges[i+1:]...)
		}
	}
}

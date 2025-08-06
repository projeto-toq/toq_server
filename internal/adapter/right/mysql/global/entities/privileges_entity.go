package globalentities

type PrivilegeEntity struct {
	ID      int64
	RoleID  int64
	Service uint8
	Method  uint8
	Allowed uint8
}

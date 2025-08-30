package usermodel

import (
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// Defines the interface methods for the userDomain interface
type UserInterface interface {
	GetID() int64
	SetID(int64)
	GetActiveRole() permissionmodel.UserRoleInterface
	SetActiveRole(active permissionmodel.UserRoleInterface)
	GetFullName() string
	SetFullName(string)
	GetNickName() string
	SetNickName(string)
	GetNationalID() string
	SetNationalID(string)
	GetCreciNumber() string
	SetCreciNumber(string)
	GetCreciState() string
	SetCreciState(string)
	GetCreciValidity() time.Time
	SetCreciValidity(time.Time)
	GetBornAt() time.Time
	SetBornAt(time.Time)
	GetPhoneNumber() string
	SetPhoneNumber(string)
	GetEmail() string
	SetEmail(string)
	GetZipCode() string
	SetZipCode(string)
	GetStreet() string
	SetStreet(string)
	GetNumber() string
	SetNumber(string)
	GetComplement() string
	SetComplement(string)
	GetNeighborhood() string
	SetNeighborhood(string)
	GetCity() string
	SetCity(string)
	GetState() string
	SetState(state string)
	GetPassword() string
	SetPassword(password string)
	IsOptStatus() bool
	SetOptStatus(bool)
	GetLastActivityAt() time.Time
	SetLastActivityAt(time.Time)
	IsDeleted() bool
	SetDeleted(bool)
	GetLastSignInAttempt() time.Time
	SetLastSignInAttempt(time.Time)
	GetDeviceToken() string
	SetDeviceToken(string)
	GetDeviceTokens() []DeviceTokenInterface
	SetDeviceTokens([]DeviceTokenInterface)
	AddDeviceToken(string) bool
}

// Creates a new UserDomain interface
func NewUser() UserInterface {
	return &user{}
}

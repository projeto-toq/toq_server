package usermodel

import (
	"time"
)

// Defines the interface methods for the userDomain interface
type UserInterface interface {
	GetID() int64
	SetID(int64)
	GetActiveRole() UserRoleInterface
	SetActiveRole(active UserRoleInterface)
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
	GetDeviceToken() string
	SetDeviceToken(string)
	GetDeviceTokens() []DeviceTokenInterface
	SetDeviceTokens([]DeviceTokenInterface)
	AddDeviceToken(string) bool

	// ==================== NEW: User-level blocking methods ====================

	// GetBlockedUntil returns the timestamp until which user is temporarily blocked
	// Returns nil if user is not temporarily blocked
	// User is blocked while GetBlockedUntil() > time.Now()
	GetBlockedUntil() *time.Time

	// SetBlockedUntil sets the temporary block expiration timestamp
	// Pass nil to clear temporary block
	// Pass &future_time to block user until that time
	SetBlockedUntil(*time.Time)

	// IsPermanentlyBlocked returns true if user is permanently blocked by admin
	// Permanent block has no expiration (requires manual admin unblock)
	IsPermanentlyBlocked() bool

	// SetPermanentlyBlocked sets or clears permanent admin block
	// true = block permanently, false = unblock
	SetPermanentlyBlocked(bool)

	// IsBlocked checks if user is currently blocked (either temporary or permanent)
	// Returns true if ANY of these conditions is true:
	//   - permanently_blocked = true
	//   - blocked_until != nil AND blocked_until > NOW()
	// This is a convenience method (can be computed from GetBlockedUntil + IsPermanentlyBlocked)
	IsBlocked() bool
}

// Creates a new UserDomain interface
func NewUser() UserInterface {
	return &user{}
}

package usermodel

import (
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

type user struct {
	id               int64
	activeRole       permissionmodel.UserRoleInterface
	fullName         string
	nickName         string
	nationalID       string
	creciNumber      string
	creciState       string
	creciValidity    time.Time
	bornAt           time.Time
	phoneNumber      string
	email            string
	zipCode          string
	street           string
	number           string
	complement       string
	neighborhood     string
	city             string
	state            string
	password         string
	optStatus        bool
	lastActivityAt   time.Time
	deleted          bool
	lastSiginAttempt time.Time
	deviceTokens     []DeviceTokenInterface
}

func (u *user) GetID() int64 {
	return u.id
}

func (u *user) SetID(id int64) {
	u.id = id
}
func (u *user) GetActiveRole() permissionmodel.UserRoleInterface {
	return u.activeRole
}

func (u *user) SetActiveRole(activeRole permissionmodel.UserRoleInterface) {
	u.activeRole = activeRole
}

func (u *user) GetFullName() string {
	return u.fullName
}

func (u *user) SetFullName(fullName string) {
	u.fullName = fullName
}

func (u *user) GetNickName() string {
	return u.nickName
}

func (u *user) SetNickName(nickName string) {
	u.nickName = nickName
}

func (u *user) GetNationalID() string {
	return u.nationalID
}

func (u *user) SetNationalID(nationalId string) {
	u.nationalID = nationalId
}

func (u *user) GetCreciNumber() string {
	return u.creciNumber
}

func (u *user) SetCreciNumber(creciNumber string) {
	u.creciNumber = creciNumber
}

func (u *user) GetCreciState() string {
	return u.creciState
}

func (u *user) SetCreciState(creciState string) {
	u.creciState = creciState
}

func (u *user) GetCreciValidity() time.Time {
	return u.creciValidity
}

func (u *user) SetCreciValidity(creciValidity time.Time) {
	u.creciValidity = creciValidity
}

func (u *user) GetBornAt() time.Time {
	return u.bornAt
}

func (u *user) SetBornAt(bornAt time.Time) {
	u.bornAt = bornAt
}

func (u *user) GetPhoneNumber() string {
	return u.phoneNumber
}

func (u *user) SetPhoneNumber(phoneNumber string) {
	u.phoneNumber = phoneNumber
}

func (u *user) GetEmail() string {
	return u.email
}

func (u *user) SetEmail(email string) {
	u.email = email
}

func (u *user) GetZipCode() string {
	return u.zipCode
}

func (u *user) SetZipCode(zipCode string) {
	u.zipCode = zipCode
}

func (u *user) GetStreet() string {
	return u.street
}

func (u *user) SetStreet(street string) {
	u.street = street
}

func (u *user) GetNumber() string {
	return u.number
}

func (u *user) SetNumber(number string) {
	u.number = number
}

func (u *user) GetComplement() string {
	return u.complement
}

func (u *user) SetComplement(complement string) {
	u.complement = complement
}

func (u *user) GetNeighborhood() string {
	return u.neighborhood
}

func (u *user) SetNeighborhood(neighborhood string) {
	u.neighborhood = neighborhood
}

func (u *user) GetCity() string {
	return u.city
}

func (u *user) SetCity(city string) {
	u.city = city
}

func (u *user) GetState() string {
	return u.state
}

func (u *user) SetState(state string) {
	u.state = state
}

func (u *user) GetPassword() string {
	return u.password
}

func (u *user) SetPassword(password string) {
	u.password = password
}
func (u *user) IsOptStatus() bool {
	return u.optStatus
}

func (u *user) SetOptStatus(optStatus bool) {
	u.optStatus = optStatus
}

func (u *user) GetLastActivityAt() time.Time {
	return u.lastActivityAt
}

func (u *user) SetLastActivityAt(lastActivityAt time.Time) {
	u.lastActivityAt = lastActivityAt
}

func (u *user) IsDeleted() bool {
	return u.deleted
}

func (u *user) SetDeleted(deleted bool) {
	u.deleted = deleted
}
func (u *user) GetLastSignInAttempt() time.Time {
	return u.lastSiginAttempt
}

func (u *user) SetLastSignInAttempt(lastSignInAttempt time.Time) {
	u.lastSiginAttempt = lastSignInAttempt
}

// Backwards-compatible single token accessor (returns first token if exists)
func (u *user) GetDeviceToken() string {
	if len(u.deviceTokens) == 0 {
		return ""
	}
	return u.deviceTokens[0].GetDeviceToken()
}

// Backwards-compatible setter: sets or replaces the first token
func (u *user) SetDeviceToken(token string) {
	if len(u.deviceTokens) == 0 {
		dt := NewDeviceToken()
		dt.SetDeviceToken(token)
		u.deviceTokens = []DeviceTokenInterface{dt}
		return
	}
	u.deviceTokens[0].SetDeviceToken(token)
}

// Full slice getter
func (u *user) GetDeviceTokens() []DeviceTokenInterface {
	return u.deviceTokens
}

// Replace all device tokens
func (u *user) SetDeviceTokens(tokens []DeviceTokenInterface) {
	u.deviceTokens = tokens
}

// AddDeviceToken adds a token if not already present; returns true if added
func (u *user) AddDeviceToken(token string) bool {
	for _, t := range u.deviceTokens {
		if t.GetDeviceToken() == token {
			return false
		}
	}
	dt := NewDeviceToken()
	dt.SetDeviceToken(token)
	u.deviceTokens = append(u.deviceTokens, dt)
	return true
}

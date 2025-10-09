package dto

// User DTOs

// CreateOwnerRequest represents owner creation request
type CreateOwnerRequest struct {
	Owner       UserCreateRequest `json:"owner" binding:"required"`
	DeviceToken string            `json:"deviceToken" binding:"required"`
}

// CreateOwnerResponse represents owner creation response
type CreateOwnerResponse struct {
	Tokens TokensResponse `json:"tokens"`
}

// CreateRealtorRequest represents realtor creation request
type CreateRealtorRequest struct {
	Realtor     UserCreateRequest `json:"realtor" binding:"required"`
	DeviceToken string            `json:"deviceToken" binding:"required"`
}

// CreateRealtorResponse represents realtor creation response
type CreateRealtorResponse struct {
	Tokens TokensResponse `json:"tokens"`
}

// CreateAgencyRequest represents agency creation request
type CreateAgencyRequest struct {
	Agency      UserCreateRequest `json:"agency" binding:"required"`
	DeviceToken string            `json:"deviceToken,omitempty"`
}

// CreateAgencyResponse represents agency creation response
type CreateAgencyResponse struct {
	Tokens TokensResponse `json:"tokens"`
}

// UserCreateRequest represents user creation data
type UserCreateRequest struct {
	FullName      string `json:"fullName,omitempty" binding:"omitempty,min=2,max=100"`
	NickName      string `json:"nickName" binding:"required,min=3,max=50"`
	NationalID    string `json:"nationalID" binding:"required"`
	CreciNumber   string `json:"creciNumber,omitempty"`
	CreciState    string `json:"creciState,omitempty"`
	CreciValidity string `json:"creciValidity,omitempty"`   // format: 2006-01-02
	BornAt        string `json:"bornAt" binding:"required"` // format: 2006-01-02
	PhoneNumber   string `json:"phoneNumber" binding:"required" example:"+5511999999999" description:"Phone number in E.164 format (e.g., +5511999999999)"`
	Email         string `json:"email" binding:"required,email"`
	ZipCode       string `json:"zipCode" binding:"required"`
	Street        string `json:"street,omitempty"`
	Number        string `json:"number,omitempty"`
	Complement    string `json:"complement,omitempty"`
	Neighborhood  string `json:"neighborhood,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty" binding:"omitempty,len=2"`
	Password      string `json:"password" binding:"required,min=8"`
}

// SignInRequest represents sign in request with user credentials
//
// Example:
//
//	{
//	  "nationalID": "12345678901",
//	  "password": "securePassword123",
//	  "deviceToken": "fcm_token_optional"
//	}
type SignInRequest struct {
	NationalID  string `json:"nationalID" binding:"required" example:"12345678901" description:"User's CPF or CNPJ (punctuation ignored; digits-only used)"`
	Password    string `json:"password" binding:"required" example:"securePassword123" description:"User's password"`
	DeviceToken string `json:"deviceToken" binding:"required" example:"fcm_device_token" description:"FCM device token for push notifications; requires X-Device-Id header"`
}

// SignInResponse represents successful sign in response with authentication tokens
//
// Example:
//
//	{
//	  "tokens": {
//	    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//	    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//	  }
//	}
type SignInResponse struct {
	Tokens TokensResponse `json:"tokens" description:"Authentication tokens"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	Tokens TokensResponse `json:"tokens"`
}

// SignOutRequest represents sign out request
type SignOutRequest struct {
	DeviceToken  string `json:"deviceToken,omitempty" example:"fcm_device_token" description:"Optional FCM device token to target a specific device session; pair with X-Device-Id header when available"`
	RefreshToken string `json:"refreshToken,omitempty" description:"Optional refresh token to revoke a single session"`
}

// SignOutResponse represents sign out response
type SignOutResponse struct {
	Message string `json:"message"`
}

// RequestPasswordChangeRequest represents password change request
type RequestPasswordChangeRequest struct {
	NationalID string `json:"nationalID" binding:"required" description:"User's CPF or CNPJ (punctuation ignored; digits-only used)"`
}

// RequestPasswordChangeResponse represents password change request response
type RequestPasswordChangeResponse struct {
	Message string `json:"message"`
}

// ConfirmPasswordChangeRequest represents password change confirmation
type ConfirmPasswordChangeRequest struct {
	NationalID  string `json:"nationalID" binding:"required" description:"User's CPF or CNPJ (punctuation ignored; digits-only used)"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
	Code        string `json:"code" binding:"required"`
}

// ConfirmPasswordChangeResponse represents password change confirmation response
type ConfirmPasswordChangeResponse struct {
	Message string `json:"message"`
}

// GetProfileResponse represents user profile response
type GetProfileResponse struct {
	User UserResponse `json:"user"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	User UpdateUserRequest `json:"user" binding:"required"`
}

// UpdateProfileResponse represents profile update response
type UpdateProfileResponse struct {
	Message string `json:"message"`
}

// UpdateUserRequest represents user update data
type UpdateUserRequest struct {
	NickName     string `json:"nickName,omitempty" binding:"omitempty,min=2,max=50"`
	BornAt       string `json:"bornAt,omitempty"` // format: 2006-01-02
	ZipCode      string `json:"zipCode,omitempty"`
	Street       string `json:"street,omitempty"`
	Number       string `json:"number,omitempty"`
	Complement   string `json:"complement,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty" binding:"omitempty,len=2"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID            int64            `json:"id"`
	ActiveRole    UserRoleResponse `json:"activeRole"`
	FullName      string           `json:"fullName"`
	NickName      string           `json:"nickName"`
	NationalID    string           `json:"nationalID"`
	CreciNumber   string           `json:"creciNumber,omitempty"`
	CreciState    string           `json:"creciState,omitempty"`
	CreciValidity string           `json:"creciValidity,omitempty"`
	BornAt        string           `json:"bornAt"`
	PhoneNumber   string           `json:"phoneNumber"`
	Email         string           `json:"email"`
	ZipCode       string           `json:"zipCode"`
	Street        string           `json:"street"`
	Number        string           `json:"number"`
	Complement    string           `json:"complement,omitempty"`
	Neighborhood  string           `json:"neighborhood"`
	City          string           `json:"city"`
	State         string           `json:"state"`
	LastActivity  string           `json:"lastActivity"`
}

// UserRoleResponse represents user role data
type UserRoleResponse struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"userId"`
	BaseRoleID   int64  `json:"baseRoleId"`
	Role         string `json:"role"`
	Active       bool   `json:"active"`
	Status       string `json:"status"`
	StatusReason string `json:"statusReason,omitempty"`
}

// DeleteAccountResponse represents account deletion response
type DeleteAccountResponse struct {
	Tokens  TokensResponse `json:"tokens"`
	Message string         `json:"message"`
}

// GetUserRolesResponse represents user roles response
type GetUserRolesResponse struct {
	Roles []UserRoleResponse `json:"roles"`
}

// GoHomeResponse represents home response
type GoHomeResponse struct {
	Message string `json:"message"`
}

// UpdateOptStatusRequest represents opt status update request
type UpdateOptStatusRequest struct {
	OptIn bool `json:"optIn"`
}

// UpdateOptStatusResponse represents opt status update response
type UpdateOptStatusResponse struct {
	Message string `json:"message"`
}

// GetPhotoUploadURLRequest represents photo upload URL request
// variant must be one of: original|small|medium|large
type GetPhotoUploadURLRequest struct {
	Variant     string `json:"variant" binding:"required,oneof=original small medium large"`
	ContentType string `json:"contentType" binding:"required"`
}

// GetPhotoUploadURLResponse represents photo upload URL response
type GetPhotoUploadURLResponse struct {
	SignedURL string `json:"signedUrl"`
}

// GetPhotoDownloadURLRequest represents photo single download URL request
// variant must be one of: original|small|medium|large
type GetPhotoDownloadURLRequest struct {
	Variant string `json:"variant" binding:"required,oneof=original small medium large"`
}

// GetPhotoDownloadURLResponse represents photo single download URL response
type GetPhotoDownloadURLResponse struct {
	SignedURL string `json:"signedUrl"`
}

// Email change requests
type RequestEmailChangeRequest struct {
	NewEmail string `json:"newEmail" binding:"required,email"`
}

type RequestEmailChangeResponse struct {
	Message string `json:"message"`
}

type ConfirmEmailChangeRequest struct {
	Code string `json:"code" binding:"required"`
}

type ConfirmEmailChangeResponse struct {
	Message string `json:"message"`
}

type ResendEmailChangeCodeResponse struct {
	Message string `json:"message"`
}

// Phone change requests
type RequestPhoneChangeRequest struct {
	NewPhoneNumber string `json:"newPhoneNumber" binding:"required" example:"+5511999999999" description:"New phone number in E.164 format (e.g., +5511999999999)"`
}

type RequestPhoneChangeResponse struct {
	Message string `json:"message"`
}

type ConfirmPhoneChangeRequest struct {
	Code string `json:"code" binding:"required"`
}

type ConfirmPhoneChangeResponse struct {
	Message string `json:"message"`
}

type ResendPhoneChangeCodeResponse struct {
	Message string `json:"message"`
}

// Role management
type AddAlternativeUserRoleRequest struct {
	CreciNumber   string `json:"creciNumber" binding:"required"`
	CreciState    string `json:"creciState" binding:"required,len=2"`
	CreciValidity string `json:"creciValidity" binding:"required"` // format: 2006-01-02
}

type AddAlternativeUserRoleResponse struct {
	Message string `json:"message"`
}

type SwitchUserRoleRequest struct {
	RoleSlug string `json:"roleSlug" binding:"required,oneof=owner realtor" example:"realtor" enum:"owner,realtor" description:"Slug do role desejado (owner ou realtor)"`
}

type SwitchUserRoleResponse struct {
	Tokens TokensResponse `json:"tokens"`
}

// Agency operations
type GetDocumentsUploadURLRequest struct {
	DocumentType string `json:"documentType" binding:"required"` // "cnpj.jpg", "front.jpg", "back.jpg", "contract.jpg"
	ContentType  string `json:"contentType" binding:"required"`  // Ex: "image/jpeg"
}

type GetDocumentsUploadURLResponse struct {
	SignedURL string `json:"signedUrl"`
}

type InviteRealtorRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required" example:"+5511999999999" description:"Realtor phone number in E.164 format (e.g., +5511999999999)"`
}

type InviteRealtorResponse struct {
	Message string `json:"message"`
}

type GetRealtorsByAgencyResponse struct {
	Realtors []UserResponse `json:"realtors"`
}

type GetRealtorByIDResponse struct {
	Realtor UserResponse `json:"realtor"`
}

type DeleteRealtorByIDResponse struct {
	Message string `json:"message"`
}

// Realtor operations
type VerifyCreciImagesResponse struct {
	Message string `json:"message"`
}

type GetCreciUploadURLRequest struct {
	DocumentType string `json:"documentType" binding:"required"` // "selfie.jpg", "front.jpg", "back.jpg"
	ContentType  string `json:"contentType" binding:"required"`  // Ex: "image/jpeg"
}

type GetCreciUploadURLResponse struct {
	SignedURL string `json:"signedUrl"`
}

type AcceptInvitationResponse struct {
	Message string `json:"message"`
}

type RejectInvitationResponse struct {
	Message string `json:"message"`
}

type GetAgencyOfRealtorResponse struct {
	Agency UserResponse `json:"agency"`
}

type DeleteAgencyOfRealtorResponse struct {
	Message string `json:"message"`
}

// User Role Status Management DTOs

// UpdateUserRoleStatusRequest represents a request to update user role status
type UpdateUserRoleStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateUserRoleStatusResponse represents response for status update
type UpdateUserRoleStatusResponse struct {
	UserID   int64  `json:"userId"`
	RoleSlug string `json:"roleSlug"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

// GetUserRoleStatusResponse represents response for status query
type GetUserRoleStatusResponse struct {
	UserID   int64  `json:"userId"`
	RoleSlug string `json:"roleSlug"`
	Status   string `json:"status"`
}

// Simple user status (active role status) for GET /user/status
// Alinhado ao handler get_user_status.go
type UserStatusResponse struct {
	Data UserStatusData `json:"data"`
}

// UserStatusData holds the minimal status payload.
// Enum:
// 0 = active
// 1 = blocked
// 2 = temp_blocked
// 3 = pending_both
// 4 = pending_email
// 5 = pending_phone
// 6 = pending_creci
// 7 = pending_cnpj
// 8 = pending_manual
// 9 = rejected
// 10 = refused_image
// 11 = refused_document
// 12 = refused_data
// 13 = deleted
type UserStatusData struct {
	Status int `json:"status" example:"0"`
}

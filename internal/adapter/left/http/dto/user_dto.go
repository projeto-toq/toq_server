package dto

// User DTOs

// CreateOwnerRequest represents owner creation request
type CreateOwnerRequest struct {
	Owner UserCreateRequest `json:"owner" binding:"required"`
}

// CreateOwnerResponse represents owner creation response
type CreateOwnerResponse struct {
	Tokens TokensResponse `json:"tokens"`
}

// CreateRealtorRequest represents realtor creation request
type CreateRealtorRequest struct {
	Realtor UserCreateRequest `json:"realtor" binding:"required"`
}

// CreateRealtorResponse represents realtor creation response
type CreateRealtorResponse struct {
	Tokens TokensResponse `json:"tokens"`
}

// CreateAgencyRequest represents agency creation request
type CreateAgencyRequest struct {
	Agency UserCreateRequest `json:"agency" binding:"required"`
}

// CreateAgencyResponse represents agency creation response
type CreateAgencyResponse struct {
	Tokens TokensResponse `json:"tokens"`
}

// UserCreateRequest represents user creation data
type UserCreateRequest struct {
	FullName      string `json:"fullName" binding:"required,min=2,max=100"`
	NickName      string `json:"nickName" binding:"required,min=2,max=50"`
	NationalID    string `json:"nationalID" binding:"required"`
	CreciNumber   string `json:"creciNumber,omitempty"`
	CreciState    string `json:"creciState,omitempty"`
	CreciValidity string `json:"creciValidity,omitempty"`   // format: 2006-01-02
	BornAt        string `json:"bornAt" binding:"required"` // format: 2006-01-02
	PhoneNumber   string `json:"phoneNumber" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	ZipCode       string `json:"zipCode" binding:"required"`
	Street        string `json:"street" binding:"required"`
	Number        string `json:"number" binding:"required"`
	Complement    string `json:"complement,omitempty"`
	Neighborhood  string `json:"neighborhood" binding:"required"`
	City          string `json:"city" binding:"required"`
	State         string `json:"state" binding:"required,len=2"`
	Password      string `json:"password" binding:"required,min=6"`
}

// SignInRequest represents sign in request
type SignInRequest struct {
	NationalID  string `json:"nationalID" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DeviceToken string `json:"deviceToken"`
}

// SignInResponse represents sign in response
type SignInResponse struct {
	Tokens TokensResponse `json:"tokens"`
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
	DeviceToken  string `json:"deviceToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

// SignOutResponse represents sign out response
type SignOutResponse struct {
	Message string `json:"message"`
}

// RequestPasswordChangeRequest represents password change request
type RequestPasswordChangeRequest struct {
	NationalID string `json:"nationalID" binding:"required"`
}

// RequestPasswordChangeResponse represents password change request response
type RequestPasswordChangeResponse struct {
	Message string `json:"message"`
}

// ConfirmPasswordChangeRequest represents password change confirmation
type ConfirmPasswordChangeRequest struct {
	NationalID  string `json:"nationalID" binding:"required"`
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

// GetOnboardingStatusResponse represents onboarding status response
type GetOnboardingStatusResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
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
type GetPhotoUploadURLRequest struct {
	ObjectName  string `json:"objectName" binding:"required"`
	ContentType string `json:"contentType" binding:"required"`
}

// GetPhotoUploadURLResponse represents photo upload URL response
type GetPhotoUploadURLResponse struct {
	SignedURL string `json:"signedUrl"`
}

// GetProfileThumbnailsResponse represents profile thumbnails response
type GetProfileThumbnailsResponse struct {
	OriginalURL string `json:"originalUrl"`
	SmallURL    string `json:"smallUrl"`
	MediumURL   string `json:"mediumUrl"`
	LargeURL    string `json:"largeUrl"`
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
	Tokens TokensResponse `json:"tokens"`
}

type ResendEmailChangeCodeResponse struct {
	Code string `json:"code"`
}

// Phone change requests
type RequestPhoneChangeRequest struct {
	NewPhoneNumber string `json:"newPhoneNumber" binding:"required"`
}

type RequestPhoneChangeResponse struct {
	Message string `json:"message"`
}

type ConfirmPhoneChangeRequest struct {
	Code string `json:"code" binding:"required"`
}

type ConfirmPhoneChangeResponse struct {
	Tokens TokensResponse `json:"tokens"`
}

type ResendPhoneChangeCodeResponse struct {
	Code string `json:"code"`
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
	RoleID int64 `json:"roleId" binding:"required"`
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
	PhoneNumber string `json:"phoneNumber" binding:"required"`
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

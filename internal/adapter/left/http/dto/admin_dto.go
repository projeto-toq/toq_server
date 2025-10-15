package dto

import (
	"fmt"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// Admin endpoints DTOs

// AdminGetPendingRealtorsRequest captures filters for GET /admin/user/pending
type AdminGetPendingRealtorsRequest struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
}

// AdminGetPendingRealtorsResponse represents GET /admin/user/pending response
type AdminGetPendingRealtorsResponse struct {
	Realtors   []AdminPendingRealtor `json:"realtors"`
	Pagination PaginationResponse    `json:"pagination"`
}

// AdminPendingRealtor minimal fields required by the spec
type AdminPendingRealtor struct {
	ID            int64  `json:"id"`
	NickName      string `json:"nickName"`
	FullName      string `json:"fullName"`
	NationalID    string `json:"nationalID"`
	CreciNumber   string `json:"creciNumber"`
	CreciValidity string `json:"creciValidity"`
	CreciState    string `json:"creciState"`
}

// AdminGetUserRequest represents POST /admin/user request
type AdminGetUserRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

// AdminGetUserResponse wraps full user data
type AdminGetUserResponse struct {
	User UserResponse `json:"user"`
}

// AdminApproveUserRequest represents POST /admin/user/approve request
// Note: status is an integer matching permission_model.UserRoleStatus
type AdminApproveUserRequest struct {
	ID     int64 `json:"id" binding:"required,min=1"`
	Status *int  `json:"status" binding:"required"`
}

// ToStatus converts the request status into a valid domain status.
func (r *AdminApproveUserRequest) ToStatus() (permissionmodel.UserRoleStatus, error) {
	if r == nil || r.Status == nil {
		return 0, fmt.Errorf("status is required")
	}

	status := permissionmodel.UserRoleStatus(*r.Status)
	if !permissionmodel.IsManualApprovalTarget(status) {
		return 0, fmt.Errorf("invalid status value: %d", *r.Status)
	}

	return status, nil
}

// AdminApproveUserResponse represents approval outcome
type AdminApproveUserResponse struct {
	Message string `json:"message"`
}

// AdminCreciDownloadURLRequest representa POST /admin/user/creci-download-url request
type AdminCreciDownloadURLRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

// AdminCreciDocumentURLs agrupa as URLs de download dos documentos CRECI
type AdminCreciDocumentURLs struct {
	Selfie string `json:"selfie"`
	Front  string `json:"front"`
	Back   string `json:"back"`
}

// AdminCreciDownloadURLResponse representa a resposta com URLs assinadas e tempo de expiração
type AdminCreciDownloadURLResponse struct {
	URLs             AdminCreciDocumentURLs `json:"urls"`
	ExpiresInMinutes int                    `json:"expiresInMinutes"`
}

// AdminListUsersRequest represents GET /admin/users filters
type AdminListUsersRequest struct {
	Page             int    `form:"page,default=1" binding:"min=1"`
	Limit            int    `form:"limit,default=20" binding:"min=1,max=100"`
	RoleName         string `form:"roleName"`
	RoleSlug         string `form:"roleSlug"`
	RoleStatus       *int   `form:"roleStatus"`
	IsSystemRole     *bool  `form:"isSystemRole"`
	FullName         string `form:"fullName"`
	CPF              string `form:"cpf"`
	Email            string `form:"email"`
	PhoneNumber      string `form:"phoneNumber"`
	Deleted          *bool  `form:"deleted"`
	IDFrom           *int64 `form:"idFrom" binding:"omitempty,min=1"`
	IDTo             *int64 `form:"idTo" binding:"omitempty,min=1"`
	BornAtFrom       string `form:"bornAtFrom"`
	BornAtTo         string `form:"bornAtTo"`
	LastActivityFrom string `form:"lastActivityFrom"`
	LastActivityTo   string `form:"lastActivityTo"`
}

// AdminListUsersResponse represents admin users listing payload
type AdminListUsersResponse struct {
	Users      []AdminUserSummary `json:"users"`
	Pagination PaginationResponse `json:"pagination"`
}

// AdminUserSummary minimal projection for admin list
type AdminUserSummary struct {
	ID          int64               `json:"id"`
	FullName    string              `json:"fullName"`
	Email       string              `json:"email"`
	PhoneNumber string              `json:"phoneNumber"`
	CPF         string              `json:"cpf"`
	Deleted     bool                `json:"deleted"`
	Role        AdminUserRoleResume `json:"role"`
}

// AdminUserRoleResume wraps active role information
type AdminUserRoleResume struct {
	UserRoleID   int64  `json:"userRoleId"`
	RoleID       int64  `json:"roleId"`
	RoleName     string `json:"roleName"`
	RoleSlug     string `json:"roleSlug"`
	IsSystemRole bool   `json:"isSystemRole"`
	Status       string `json:"status"`
	IsActive     bool   `json:"isActive"`
}

// AdminCreateSystemUserRequest represents POST /admin/system-users request
type AdminCreateSystemUserRequest struct {
	FullName    string `json:"fullName" binding:"required,min=2,max=150"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	CPF         string `json:"cpf" binding:"required"`
	RoleSlug    string `json:"roleSlug" binding:"required"`
	BornAt      string `json:"bornAt" binding:"required"`
}

// AdminUpdateSystemUserRequest represents PUT /admin/system-users request body
type AdminUpdateSystemUserRequest struct {
	UserID      int64  `json:"userId" binding:"required,min=1"`
	FullName    string `json:"fullName" binding:"required,min=2,max=150"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

// AdminSystemUserResponse basic response for create/update actions
type AdminSystemUserResponse struct {
	UserID  int64  `json:"userId"`
	Slug    string `json:"roleSlug"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// AdminDeleteSystemUserRequest represents DELETE /admin/system-users payload
type AdminDeleteSystemUserRequest struct {
	UserID int64 `json:"userId" binding:"required,min=1"`
}

// AdminListRolesRequest represents GET /admin/roles filters
type AdminListRolesRequest struct {
	Page         int    `form:"page,default=1" binding:"min=1"`
	Limit        int    `form:"limit,default=20" binding:"min=1,max=100"`
	Name         string `form:"name"`
	Slug         string `form:"slug"`
	Description  string `form:"description"`
	IsSystemRole *bool  `form:"isSystemRole"`
	IsActive     *bool  `form:"isActive"`
	IDFrom       *int64 `form:"idFrom" binding:"omitempty,min=1"`
	IDTo         *int64 `form:"idTo" binding:"omitempty,min=1"`
}

// AdminRoleSummary minimal role representation
type AdminRoleSummary struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	IsSystemRole bool   `json:"isSystemRole"`
	IsActive     bool   `json:"isActive"`
}

// AdminListRolesResponse wraps role listing
type AdminListRolesResponse struct {
	Roles      []AdminRoleSummary `json:"roles"`
	Pagination PaginationResponse `json:"pagination"`
}

// AdminCreateRoleRequest represents POST /admin/roles request body
type AdminCreateRoleRequest struct {
	Name         string `json:"name" binding:"required,min=2,max=100"`
	Slug         string `json:"slug" binding:"required"`
	Description  string `json:"description"`
	IsSystemRole bool   `json:"isSystemRole"`
}

// AdminUpdateRoleRequest represents PUT /admin/roles request body
type AdminUpdateRoleRequest struct {
	ID          int64  `json:"id" binding:"required,min=1"`
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

// AdminDeleteRoleRequest represents DELETE /admin/roles request body
type AdminDeleteRoleRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

// AdminRestoreRoleRequest represents POST /admin/roles/restore request body
type AdminRestoreRoleRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

// AdminRoleResponse basic role response payload
type AdminRoleResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

// AdminListPermissionsRequest captures filters for GET /admin/permissions
type AdminListPermissionsRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	Limit    int    `form:"limit,default=20" binding:"min=1,max=100"`
	Name     string `form:"name"`
	Resource string `form:"resource"`
	Action   string `form:"action"`
	IsActive *bool  `form:"isActive"`
}

// AdminPermissionSummary minimal projection for permission
type AdminPermissionSummary struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Resource    string         `json:"resource"`
	Action      string         `json:"action"`
	Description string         `json:"description"`
	Conditions  map[string]any `json:"conditions,omitempty"`
	IsActive    bool           `json:"isActive"`
}

// AdminListPermissionsResponse wraps permission listing
type AdminListPermissionsResponse struct {
	Permissions []AdminPermissionSummary `json:"permissions"`
	Pagination  PaginationResponse       `json:"pagination"`
}

// AdminCreatePermissionRequest represents POST /admin/permissions payload
type AdminCreatePermissionRequest struct {
	Name        string         `json:"name" binding:"required,min=2,max=100"`
	Resource    string         `json:"resource" binding:"required,min=1,max=50"`
	Action      string         `json:"action" binding:"required,min=1,max=50"`
	Description string         `json:"description"`
	Conditions  map[string]any `json:"conditions"`
}

// AdminUpdatePermissionRequest represents PUT /admin/permissions payload
type AdminUpdatePermissionRequest struct {
	ID          int64          `json:"id" binding:"required,min=1"`
	Name        string         `json:"name" binding:"required,min=2,max=100"`
	Description string         `json:"description"`
	IsActive    *bool          `json:"isActive"`
	Conditions  map[string]any `json:"conditions"`
}

// AdminDeletePermissionRequest represents DELETE /admin/permissions payload
type AdminDeletePermissionRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

// AdminPermissionResponse basic response for permission operations
type AdminPermissionResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

// AdminListRolePermissionsRequest filters GET /admin/role-permissions
type AdminListRolePermissionsRequest struct {
	Page         int    `form:"page,default=1" binding:"min=1"`
	Limit        int    `form:"limit,default=20" binding:"min=1,max=100"`
	RoleID       *int64 `form:"roleId" binding:"omitempty,min=1"`
	PermissionID *int64 `form:"permissionId" binding:"omitempty,min=1"`
	Granted      *bool  `form:"granted"`
}

// AdminRolePermissionSummary minimal projection for role-permission links
type AdminRolePermissionSummary struct {
	ID           int64          `json:"id"`
	RoleID       int64          `json:"roleId"`
	PermissionID int64          `json:"permissionId"`
	Granted      bool           `json:"granted"`
	Conditions   map[string]any `json:"conditions,omitempty"`
}

// AdminListRolePermissionsResponse wraps role-permission listing
type AdminListRolePermissionsResponse struct {
	RolePermissions []AdminRolePermissionSummary `json:"rolePermissions"`
	Pagination      PaginationResponse           `json:"pagination"`
}

// AdminCreateRolePermissionRequest represents POST /admin/role-permissions payload
type AdminCreateRolePermissionRequest struct {
	RoleID       int64          `json:"roleId" binding:"required,min=1"`
	PermissionID int64          `json:"permissionId" binding:"required,min=1"`
	Granted      *bool          `json:"granted"`
	Conditions   map[string]any `json:"conditions"`
}

// AdminUpdateRolePermissionRequest represents PUT /admin/role-permissions payload
type AdminUpdateRolePermissionRequest struct {
	ID         int64          `json:"id" binding:"required,min=1"`
	Granted    *bool          `json:"granted"`
	Conditions map[string]any `json:"conditions"`
}

// AdminDeleteRolePermissionRequest represents DELETE /admin/role-permissions payload
type AdminDeleteRolePermissionRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

// AdminRolePermissionResponse basic response for role-permission operations
type AdminRolePermissionResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

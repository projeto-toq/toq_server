package dto

import (
	"fmt"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// Admin endpoints DTOs

// AdminGetPendingRealtorsResponse represents GET /admin/user/pending response
type AdminGetPendingRealtorsResponse struct {
	Realtors []AdminPendingRealtor `json:"realtors"`
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

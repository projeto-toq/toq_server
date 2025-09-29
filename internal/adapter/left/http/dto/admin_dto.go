package dto

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
	Status int   `json:"status" binding:"required"`
}

// AdminApproveUserResponse represents approval outcome
type AdminApproveUserResponse struct {
	Message string `json:"message"`
}

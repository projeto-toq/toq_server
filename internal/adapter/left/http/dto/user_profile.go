package dto

// UserProfileResponse represents the response envelope for the user profile endpoint.
// It contains a single Data field with the user profile details.
type UserProfileResponse struct {
	Data UserProfileData `json:"data"`
}

// UserProfileData holds the user profile attributes returned to clients.
// Use of concrete types and json tags ensures a stable contract.
type UserProfileData struct {
	ID          int64          `json:"id"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phone_number"`
	FullName    string         `json:"full_name"`
	NickName    string         `json:"nick_name"`
	NationalID  string         `json:"national_id"`
	ActiveRole  *ActiveRoleDTO `json:"active_role,omitempty"`
	BornAt      string         `json:"born_at,omitempty"`
	ZipCode     string         `json:"zip_code"`
	Street      string         `json:"street"`
	City        string         `json:"city"`
	State       string         `json:"state"`
}

// ActiveRoleDTO represents the user's active role in a safe, client-facing form.
type ActiveRoleDTO struct {
	ID     int64    `json:"id"`
	Role   *RoleDTO `json:"role,omitempty"`
	Active bool     `json:"active"`
}

// RoleDTO represents a role with key attributes exposed to clients.
type RoleDTO struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description,omitempty"`
	IsSystemRole bool   `json:"is_system_role"`
	IsActive     bool   `json:"is_active"`
}

// UserStatusResponse represents the response for GET /user/status endpoint.
// It returns only the current active role status of the authenticated user.
type UserStatusResponse struct {
	Data UserStatusData `json:"data"`
}

// UserStatusData holds the minimal status payload.
type UserStatusData struct {
	// Status código numérico do status da role ativa.
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
	// 9 = deleted
	// 10 = invite_pending
	Status int `json:"status" example:"0"`
}

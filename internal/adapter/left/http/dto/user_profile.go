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

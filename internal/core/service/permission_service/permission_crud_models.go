package permissionservice

import permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"

// ListPermissionsInput aggregates filters for permission pagination.
type ListPermissionsInput struct {
	Page     int
	Limit    int
	Name     string
	Action   string
	IsActive *bool
}

// ListPermissionsOutput provides paginated permission results.
type ListPermissionsOutput struct {
	Permissions []permissionmodel.PermissionInterface
	Total       int64
	Page        int
	Limit       int
}

// CreatePermissionInput carries data to create a new permission.
type CreatePermissionInput struct {
	Name        string
	Action      string
	Description string
}

// UpdatePermissionInput carries data to update an existing permission.
type UpdatePermissionInput struct {
	ID          int64
	Name        string
	Description string
	IsActive    *bool
}

// ListRolePermissionsInput aggregates filters for role-permission pagination.
type ListRolePermissionsInput struct {
	Page         int
	Limit        int
	RoleID       *int64
	PermissionID *int64
	Granted      *bool
}

// ListRolePermissionsOutput provides paginated role-permission results.
type ListRolePermissionsOutput struct {
	RolePermissions []permissionmodel.RolePermissionInterface
	Total           int64
	Page            int
	Limit           int
}

// CreateRolePermissionInput carries data to create a role-permission relation.
type CreateRolePermissionInput struct {
	RoleID       int64
	PermissionID int64
	Granted      *bool
	Conditions   map[string]any
}

// UpdateRolePermissionInput carries data to update a role-permission relation.
type UpdateRolePermissionInput struct {
	ID         int64
	Granted    *bool
	Conditions map[string]any
}

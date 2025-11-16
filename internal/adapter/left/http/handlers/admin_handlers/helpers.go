package adminhandlers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

func computeTotalPages(total int64, limit int) int {
	if limit <= 0 || total <= 0 {
		return 0
	}

	pages := int(total / int64(limit))
	if total%int64(limit) != 0 {
		pages++
	}

	if pages == 0 && total > 0 {
		return 1
	}

	return pages
}

func toAdminUserSummary(user usermodel.UserInterface) dto.AdminUserSummary {
	summary := dto.AdminUserSummary{
		ID:          user.GetID(),
		FullName:    user.GetFullName(),
		Email:       user.GetEmail(),
		PhoneNumber: user.GetPhoneNumber(),
		CPF:         user.GetNationalID(),
		Deleted:     user.IsDeleted(),
	}

	if active := user.GetActiveRole(); active != nil {
		roleResume := dto.AdminUserRoleResume{
			UserRoleID: active.GetID(),
			RoleID:     active.GetRoleID(),
			Status:     active.GetStatus().String(),
			IsActive:   active.GetIsActive(),
		}
		if role := active.GetRole(); role != nil {
			roleResume.RoleName = role.GetName()
			roleResume.RoleSlug = role.GetSlug()
			roleResume.IsSystemRole = role.GetIsSystemRole()
		}
		summary.Role = roleResume
	}

	return summary
}

func toAdminRoleSummary(role permissionmodel.RoleInterface) dto.AdminRoleSummary {
	if role == nil {
		return dto.AdminRoleSummary{}
	}

	return dto.AdminRoleSummary{
		ID:           role.GetID(),
		Name:         role.GetName(),
		Slug:         role.GetSlug(),
		Description:  role.GetDescription(),
		IsSystemRole: role.GetIsSystemRole(),
		IsActive:     role.GetIsActive(),
	}
}

func toAdminPermissionSummary(permission permissionmodel.PermissionInterface) dto.AdminPermissionSummary {
	if permission == nil {
		return dto.AdminPermissionSummary{}
	}

	resp := dto.AdminPermissionSummary{
		ID:          permission.GetID(),
		Name:        permission.GetName(),
		Action:      permission.GetAction(),
		Description: permission.GetDescription(),
		IsActive:    permission.GetIsActive(),
	}

	return resp
}

func toAdminRolePermissionSummary(rolePermission permissionmodel.RolePermissionInterface) dto.AdminRolePermissionSummary {
	if rolePermission == nil {
		return dto.AdminRolePermissionSummary{}
	}

	resp := dto.AdminRolePermissionSummary{
		ID:           rolePermission.GetID(),
		RoleID:       rolePermission.GetRoleID(),
		PermissionID: rolePermission.GetPermissionID(),
		Granted:      rolePermission.GetGranted(),
	}

	return resp
}

// filterRoutes applies method and path pattern filters to a list of Gin routes.
// Both filters are case-insensitive. If both filters are empty, returns all routes.
//
// Parameters:
//   - routes: Full list of routes from Gin engine
//   - method: HTTP method filter (empty string = no filter)
//   - pathPattern: Path substring filter (empty string = no filter)
//
// Returns:
//   - Filtered list of routes matching the criteria
func filterRoutes(routes gin.RoutesInfo, method, pathPattern string) gin.RoutesInfo {
	// No filters: return all routes
	if method == "" && pathPattern == "" {
		return routes
	}

	filtered := make(gin.RoutesInfo, 0, len(routes))
	for _, route := range routes {
		// Method filter (case-insensitive exact match)
		if method != "" && !strings.EqualFold(route.Method, method) {
			continue
		}

		// Path pattern filter (case-insensitive substring match)
		if pathPattern != "" && !strings.Contains(strings.ToLower(route.Path), strings.ToLower(pathPattern)) {
			continue
		}

		filtered = append(filtered, route)
	}
	return filtered
}

// extractHandlerName cleans up the Gin handler function name for display.
// Converts full package path to concise "Struct.Method" format.
//
// Input example:
//
//	"github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/admin_handlers.(*AdminHandler).GetAdminPermissions-fm"
//
// Output example:
//
//	"AdminHandler.GetAdminPermissions"
//
// If parsing fails, returns the original handler name as fallback.
//
// Parameters:
//   - handlerName: Full handler function name from Gin RouteInfo
//
// Returns:
//   - Cleaned handler name (Struct.Method) or original if parsing fails
func extractHandlerName(handlerName string) string {
	// Split by dots to get package components
	parts := strings.Split(handlerName, ".")
	if len(parts) < 2 {
		// Fallback: return original if format is unexpected
		return handlerName
	}

	// Get last two parts: (*Struct) and Method-fm
	structPart := parts[len(parts)-2]
	methodPart := parts[len(parts)-1]

	// Clean struct name: (*AdminHandler) -> AdminHandler
	structPart = strings.Trim(structPart, "(*)")

	// Clean method name: GetAdminPermissions-fm -> GetAdminPermissions
	methodPart = strings.Split(methodPart, "-")[0]

	return fmt.Sprintf("%s.%s", structPart, methodPart)
}

// calculatePaginationBounds computes slice indices for pagination.
// Handles edge cases: invalid page/limit, out-of-bounds indices.
//
// Parameters:
//   - page: Page number (1-indexed, defaults to 1 if <= 0)
//   - limit: Items per page (defaults to 50 if <= 0)
//   - total: Total number of items
//
// Returns:
//   - start: Starting index for slice (inclusive)
//   - end: Ending index for slice (exclusive)
func calculatePaginationBounds(page, limit int, total int64) (start, end int) {
	// Normalize page and limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}

	// Calculate bounds
	start = (page - 1) * limit
	end = start + limit

	// Clamp to total length
	if start > int(total) {
		start = int(total)
	}
	if end > int(total) {
		end = int(total)
	}

	return start, end
}

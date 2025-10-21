package adminhandlers

import (
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

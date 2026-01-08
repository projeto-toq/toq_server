package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	permissionservice "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetAdminRolePermissions handles GET /admin/role-permissions
//
//	@Summary	List role-permission associations
//	@Tags		Admin Permissions
//	@Produce	json
//	@Param		page		query	int	false	"Page number" default(1) Extensions(x-example=1)
//	@Param		limit		query	int	false	"Page size" default(20) Extensions(x-example=20)
//	@Param		roleId		query	int	false	"Filter by role ID"
//	@Param		permissionId	query	int	false	"Filter by permission ID"
//	@Param		granted		query	bool	false	"Filter by granted flag"
//	@Success	200	{object}	dto.AdminListRolePermissionsResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/role-permissions [get]
func (h *AdminHandler) GetAdminRolePermissions(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminListRolePermissionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := permissionservice.ListRolePermissionsInput{
		Page:    req.Page,
		Limit:   req.Limit,
		Granted: req.Granted,
	}
	if req.RoleID != nil {
		input.RoleID = req.RoleID
	}
	if req.PermissionID != nil {
		input.PermissionID = req.PermissionID
	}

	result, err := h.permissionService.ListRolePermissions(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminListRolePermissionsResponse{
		RolePermissions: make([]dto.AdminRolePermissionSummary, 0, len(result.RolePermissions)),
		Pagination: dto.PaginationResponse{
			Page:       result.Page,
			Limit:      result.Limit,
			Total:      result.Total,
			TotalPages: computeTotalPages(result.Total, result.Limit),
		},
	}

	for _, rolePermission := range result.RolePermissions {
		resp.RolePermissions = append(resp.RolePermissions, toAdminRolePermissionSummary(rolePermission))
	}

	c.JSON(http.StatusOK, resp)
}

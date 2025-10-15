package adminhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	permissionservice "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetAdminPermissions handles GET /admin/permissions
//
//	@Summary	List permissions with pagination
//	@Tags		Admin
//	@Produce	json
//	@Param		page		query	int	false	"Page number" default(1) example(1)
//	@Param		limit		query	int	false	"Page size" default(20) example(20)
//	@Param		name		query	string	false	"Filter by name (supports '*' wildcard)"
//	@Param		resource	query	string	false	"Filter by resource (supports '*' wildcard)"
//	@Param		action		query	string	false	"Filter by action (supports '*' wildcard)"
//	@Param		isActive	query	bool	false	"Filter by active flag"
//	@Success	200	{object}	dto.AdminListPermissionsResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/permissions [get]
func (h *AdminHandler) GetAdminPermissions(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminListPermissionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := permissionservice.ListPermissionsInput{
		Page:     req.Page,
		Limit:    req.Limit,
		Name:     strings.TrimSpace(req.Name),
		Resource: strings.TrimSpace(req.Resource),
		Action:   strings.TrimSpace(req.Action),
		IsActive: req.IsActive,
	}

	result, err := h.permissionService.ListPermissions(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminListPermissionsResponse{
		Permissions: make([]dto.AdminPermissionSummary, 0, len(result.Permissions)),
		Pagination: dto.PaginationResponse{
			Page:       result.Page,
			Limit:      result.Limit,
			Total:      result.Total,
			TotalPages: computeTotalPages(result.Total, result.Limit),
		},
	}

	for _, permission := range result.Permissions {
		resp.Permissions = append(resp.Permissions, toAdminPermissionSummary(permission))
	}

	c.JSON(http.StatusOK, resp)
}

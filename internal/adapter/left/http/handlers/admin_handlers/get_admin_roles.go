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

// GetAdminRoles handles GET /admin/roles
//
//	@Summary      List roles with pagination
//	@Tags         Admin Roles
//	@Produce      json
//	@Param        page          query  int    false "Page number" default(1) Extensions(x-example=1)
//	@Param        limit         query  int    false "Page size" default(20) Extensions(x-example=20)
//	@Param        name          query  string false "Filter by name (supports '*' wildcard)" Extensions(x-example="*Manager*")
//	@Param        slug          query  string false "Filter by slug (supports '*' wildcard)" Extensions(x-example="*admin*")
//	@Param        description   query  string false "Filter by description (supports '*' wildcard)" Extensions(x-example="*acesso*")
//	@Param        isSystemRole  query  bool   false "Filter by system role flag" Extensions(x-example=true)
//	@Param        isActive      query  bool   false "Filter by active flag" Extensions(x-example=true)
//	@Param        idFrom        query  int    false "Filter by role ID (from)" Extensions(x-example=10)
//	@Param        idTo          query  int    false "Filter by role ID (to)" Extensions(x-example=200)
//	@Success      200  {object}  dto.AdminListRolesResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/roles [get]
func (h *AdminHandler) GetAdminRoles(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminListRolesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := permissionservice.ListRolesInput{
		Page:         req.Page,
		Limit:        req.Limit,
		Name:         strings.TrimSpace(req.Name),
		Slug:         strings.TrimSpace(req.Slug),
		Description:  strings.TrimSpace(req.Description),
		IsSystemRole: req.IsSystemRole,
		IsActive:     req.IsActive,
		IDFrom:       req.IDFrom,
		IDTo:         req.IDTo,
	}

	result, err := h.permissionService.ListRoles(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminListRolesResponse{
		Roles: make([]dto.AdminRoleSummary, 0, len(result.Roles)),
		Pagination: dto.PaginationResponse{
			Page:       input.Page,
			Limit:      input.Limit,
			Total:      result.Total,
			TotalPages: computeTotalPages(result.Total, input.Limit),
		},
	}

	for _, role := range result.Roles {
		resp.Roles = append(resp.Roles, toAdminRoleSummary(role))
	}

	c.JSON(http.StatusOK, resp)
}

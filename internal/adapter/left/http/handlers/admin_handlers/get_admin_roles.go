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
//	@Tags         Admin
//	@Produce      json
//	@Param        page          query  int    false  "Page number" default(1)
//	@Param        limit         query  int    false  "Page size" default(20)
//	@Param        name          query  string false "Filter by name"
//	@Param        slug          query  string false "Filter by slug"
//	@Param        isSystemRole  query  bool   false "Filter by system role flag"
//	@Param        isActive      query  bool   false "Filter by active flag"
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
		IsSystemRole: req.IsSystemRole,
		IsActive:     req.IsActive,
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

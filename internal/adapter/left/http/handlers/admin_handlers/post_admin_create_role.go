package adminhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostAdminCreateRole handles POST /admin/roles
//
//	@Summary      Create a role
//	@Tags         Admin
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminCreateRoleRequest  true  "Role payload"
//	@Success      201  {object}  dto.AdminRoleResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      409  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/roles [post]
func (h *AdminHandler) PostAdminCreateRole(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminCreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	name := strings.TrimSpace(req.Name)
	slug := permissionmodel.RoleSlug(strings.TrimSpace(req.Slug))
	role, err := h.permissionService.CreateRole(ctx, name, slug, req.Description, req.IsSystemRole)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminRoleResponse{
		ID:      role.GetID(),
		Message: "Role created",
	}
	c.JSON(http.StatusCreated, resp)
}

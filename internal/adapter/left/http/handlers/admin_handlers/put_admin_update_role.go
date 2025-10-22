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

// PutAdminUpdateRole handles PUT /admin/roles
//
//	@Summary      Update a role
//	@Tags         Admin Roles
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminUpdateRoleRequest  true  "Role update payload"
//	@Success      200  {object}  dto.AdminRoleResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      404  {object}  map[string]any
//	@Failure      409  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/roles [put]
func (h *AdminHandler) PutAdminUpdateRole(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminUpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := permissionservice.UpdateRoleInput{
		ID:          req.ID,
		Name:        strings.TrimSpace(req.Name),
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	role, err := h.permissionService.UpdateRole(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminRoleResponse{
		ID:      role.GetID(),
		Message: "Role updated",
	}
	c.JSON(http.StatusOK, resp)
}

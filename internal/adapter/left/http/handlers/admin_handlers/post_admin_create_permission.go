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

// PostAdminCreatePermission handles POST /admin/permissions
//
//	@Summary	Create a new permission
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		payload	body	dto.AdminCreatePermissionRequest	true	"Permission payload"
//	@Success	201	{object}	dto.AdminPermissionResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	409	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/permissions [post]
func (h *AdminHandler) PostAdminCreatePermission(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminCreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := permissionservice.CreatePermissionInput{
		Name:        strings.TrimSpace(req.Name),
		Action:      strings.TrimSpace(req.Action),
		Description: strings.TrimSpace(req.Description),
	}

	permission, err := h.permissionService.CreatePermission(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminPermissionResponse{
		ID:      permission.GetID(),
		Message: "permission created",
	}

	c.JSON(http.StatusCreated, resp)
}

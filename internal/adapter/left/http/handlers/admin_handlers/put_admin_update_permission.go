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

// PutAdminUpdatePermission handles PUT /admin/permissions
//
//	@Summary	Update an existing permission
//	@Tags		Admin
//	@Accept		json
//	@Produce	json
//	@Param		payload	body	dto.AdminUpdatePermissionRequest	true	"Permission payload"
//	@Success	200	{object}	dto.AdminPermissionResponse
//	@Failure	400	{object}	map[string]any
//	@Failure	401	{object}	map[string]any
//	@Failure	403	{object}	map[string]any
//	@Failure	404	{object}	map[string]any
//	@Failure	409	{object}	map[string]any
//	@Failure	500	{object}	map[string]any
//	@Router		/admin/permissions [put]
func (h *AdminHandler) PutAdminUpdatePermission(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.AdminUpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := permissionservice.UpdatePermissionInput{
		ID:          req.ID,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		IsActive:    req.IsActive,
		Conditions:  req.Conditions,
	}

	permission, err := h.permissionService.UpdatePermission(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	resp := dto.AdminPermissionResponse{
		ID:      permission.GetID(),
		Message: "permission updated",
	}

	c.JSON(http.StatusOK, resp)
}

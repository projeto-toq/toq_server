package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteAdminSystemUser handles DELETE /admin/users/system
//
//	@Summary      Deactivate a system user
//	@Tags         Admin
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminDeleteSystemUserRequest  true  "Deletion payload"
//	@Success      204  "System user deactivated"
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      404  {object}  map[string]any
//	@Failure      409  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/users/system [delete]
func (h *AdminHandler) DeleteAdminSystemUser(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminDeleteSystemUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}
	if req.UserID <= 0 {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("userId", "Invalid user id"))
		return
	}

	if err := h.userService.DeleteSystemUser(ctx, userservices.DeleteSystemUserInput{UserID: req.UserID}); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

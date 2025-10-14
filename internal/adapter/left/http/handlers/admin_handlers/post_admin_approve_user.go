package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostAdminApproveUser handles POST /admin/user/approve
//
//	@Summary      Approve or refuse realtor status manually
//	@Description  Status must be one of the allowed enum values. On success, sends FCM notification.
//	@Tags         Admin
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminApproveUserRequest  true  "User ID and target status (enum: 0=active,10=refused_image,11=refused_document,12=refused_data)"
//	@Success      200  {object}  dto.AdminApproveUserResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      409  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/user/approve [post]
func (h *AdminHandler) PostAdminApproveUser(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminApproveUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	target, statusErr := req.ToStatus()
	if statusErr != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("status", statusErr.Error()))
		return
	}

	if err := h.userService.ApproveCreciManual(ctx, req.ID, target); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.AdminApproveUserResponse{Message: "Status updated"})
}

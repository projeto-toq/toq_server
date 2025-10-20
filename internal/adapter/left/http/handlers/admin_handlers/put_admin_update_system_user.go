package adminhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PutAdminUpdateSystemUser handles PUT /admin/users/system
//
//	@Summary      Update system user data
//	@Tags         Admin
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminUpdateSystemUserRequest  true  "Update payload"
//	@Success      200  {object}  dto.AdminSystemUserResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      404  {object}  map[string]any
//	@Failure      409  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/users/system [put]
func (h *AdminHandler) PutAdminUpdateSystemUser(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminUpdateSystemUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	input := userservices.UpdateSystemUserInput{
		UserID:      req.UserID,
		FullName:    strings.TrimSpace(req.FullName),
		Email:       strings.TrimSpace(req.Email),
		PhoneNumber: strings.TrimSpace(req.PhoneNumber),
	}

	result, svcErr := h.userService.UpdateSystemUser(ctx, input)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	resp := dto.AdminSystemUserResponse{
		UserID:  result.UserID,
		Slug:    result.RoleSlug.String(),
		Email:   result.Email,
		Message: "System user updated",
	}
	c.JSON(http.StatusOK, resp)
}

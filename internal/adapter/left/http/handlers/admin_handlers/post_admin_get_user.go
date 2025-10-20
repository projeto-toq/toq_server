package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpconv "github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostAdminGetUser handles POST /admin/users/detail
//
//	@Summary      Get full user by ID
//	@Tags         Admin
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminGetUserRequest  true  "User ID"
//	@Success      200  {object}  dto.AdminGetUserResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      404  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/users/detail [post]
func (h *AdminHandler) PostAdminGetUser(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminGetUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	user, err := h.userService.GetUserByID(ctx, req.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	profile := httpconv.ToGetProfileResponse(user)
	c.JSON(http.StatusOK, dto.AdminGetUserResponse(profile))
}

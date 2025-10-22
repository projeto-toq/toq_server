package adminhandlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PostAdminCreateSystemUser handles POST /admin/users/system
//
//	@Summary      Create a new system user
//	@Tags         Admin Users
//	@Accept       json
//	@Produce      json
//	@Param        request  body  dto.AdminCreateSystemUserRequest  true  "System user payload"
//	@Description	Create a System User with roleSlug: (photographer, attendantRealtor, attendantOwner, attendant, manager) and details. Not for Owner/Realtor user creation. Email with instruction will be sent to the new user.
//	@Success      201  {object}  dto.AdminSystemUserResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      409  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /admin/users/system [post]
func (h *AdminHandler) PostAdminCreateSystemUser(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	var req dto.AdminCreateSystemUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	bornAt, err := time.Parse("2006-01-02", req.BornAt)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("bornAt", "Date must be in format YYYY-MM-DD"))
		return
	}

	input := userservices.CreateSystemUserInput{
		FullName:    strings.TrimSpace(req.FullName),
		Email:       strings.TrimSpace(req.Email),
		PhoneNumber: strings.TrimSpace(req.PhoneNumber),
		CPF:         strings.TrimSpace(req.CPF),
		BornAt:      bornAt,
		RoleSlug:    permissionmodel.RoleSlug(strings.TrimSpace(req.RoleSlug)),
	}

	result, svcErr := h.userService.CreateSystemUser(ctx, input)
	if svcErr != nil {
		httperrors.SendHTTPErrorObj(c, svcErr)
		return
	}

	resp := dto.AdminSystemUserResponse{
		UserID:  result.UserID,
		Slug:    result.RoleSlug.String(),
		Email:   result.Email,
		Message: "System user created",
	}
	c.JSON(http.StatusCreated, resp)
}

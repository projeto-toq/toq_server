package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) SwitchUserRole(c *gin.Context) {
	// Enrich context with request info and user
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Parse request body
	var request dto.SwitchUserRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service to switch user role
	tokens, err := uh.userService.SwitchUserRole(ctx, permissionmodel.RoleSlug(request.RoleSlug))
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.SwitchUserRoleResponse{
		Tokens: dto.TokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}

	c.JSON(http.StatusOK, response)
}

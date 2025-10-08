package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetUserRoles(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Get user information from Gin context (set by AuthMiddleware)
	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	// Call permission service directly to get user roles (no business logic required)
	roles, err := uh.permissionService.GetUserRoles(ctx, userInfo.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Convert roles to DTO format
	var roleResponses []dto.UserRoleResponse
	for _, role := range roles {
		roleResponse := dto.UserRoleResponse{
			ID:           role.GetID(),
			UserID:       role.GetUserID(),
			BaseRoleID:   role.GetRoleID(),
			Role:         role.GetRole().GetSlug(),
			Active:       role.GetIsActive(),
			Status:       role.GetStatus().String(),
			StatusReason: "", // Removed as per requirements - status is self-explanatory
		}
		roleResponses = append(roleResponses, roleResponse)
	}

	// Prepare response
	response := dto.GetUserRolesResponse{
		Roles: roleResponses,
	}

	c.JSON(http.StatusOK, response)
}

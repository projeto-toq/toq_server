package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetUserRoles(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call permission service directly to get user roles (no business logic required)
	roles, err := uh.permissionService.GetUserRoles(ctx, userInfo.ID)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "GET_USER_ROLES_FAILED", "Failed to get user roles")
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
			Status:       "active", // TODO: Implement status system
			StatusReason: "",       // TODO: Implement status reason system
		}
		roleResponses = append(roleResponses, roleResponse)
	}

	// Prepare response
	response := dto.GetUserRolesResponse{
		Roles: roleResponses,
	}

	c.JSON(http.StatusOK, response)
}

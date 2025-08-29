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

	// Call service to get user roles
	roles, err := uh.userService.GetUserRolesByUser(ctx, userInfo.ID)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "GET_USER_ROLES_FAILED", "Failed to get user roles")
		return
	}

	// Convert roles to DTO
	rolesResponse := make([]dto.UserRoleResponse, len(roles))
	for i, role := range roles {
		rolesResponse[i] = dto.UserRoleResponse{
			ID:           role.GetID(),
			UserID:       role.GetUserID(),
			BaseRoleID:   role.GetBaseRoleID(),
			Role:         role.GetRole().String(),
			Active:       role.IsActive(),
			Status:       role.GetStatus().String(),
			StatusReason: role.GetStatusReason(),
		}
	}

	// Prepare response
	response := dto.GetUserRolesResponse{
		Roles: rolesResponse,
	}

	c.JSON(http.StatusOK, response)
}

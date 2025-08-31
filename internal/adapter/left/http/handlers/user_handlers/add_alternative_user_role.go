package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) AddAlternativeUserRole(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Parse request body
	var request dto.AddAlternativeUserRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Get user to determine current role - precisamos buscar o usuário para obter o role
	user, err := uh.userService.GetProfile(ctx, userInfo.ID)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "GET_USER_FAILED", "Failed to get user information")
		return
	}

	// Determine alternative role based on current active role
	var alternativeRole permissionmodel.RoleSlug
	currentRole := utils.GetUserRoleSlugFromUserRole(user.GetActiveRole())
	if currentRole == permissionmodel.RoleSlugOwner {
		alternativeRole = permissionmodel.RoleSlugRealtor
	} else {
		alternativeRole = permissionmodel.RoleSlugOwner
	}

	// Call service to add alternative role
	if err := uh.userService.AddAlternativeRole(ctx, userInfo.ID, alternativeRole, request.CreciNumber, request.CreciState, request.CreciValidity); err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "ADD_ALTERNATIVE_ROLE_FAILED", "Failed to add alternative role")
		return
	}

	// Prepare response
	response := dto.AddAlternativeUserRoleResponse{
		Message: "Alternative user role added successfully",
	}

	c.JSON(http.StatusOK, response)
}

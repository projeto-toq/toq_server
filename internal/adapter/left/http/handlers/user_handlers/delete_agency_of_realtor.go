package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) DeleteAgencyOfRealtor(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to delete agency of realtor
	if err := uh.userService.DeleteAgencyOfRealtor(ctx, userInfo.ID); err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "DELETE_AGENCY_OF_REALTOR_FAILED", "Failed to delete agency of realtor")
		return
	}

	// Prepare response
	response := dto.DeleteAgencyOfRealtorResponse{
		Message: "Agency of realtor deleted successfully",
	}

	c.JSON(http.StatusOK, response)
}

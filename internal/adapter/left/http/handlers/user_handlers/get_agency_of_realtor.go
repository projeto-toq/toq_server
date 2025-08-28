package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetAgencyOfRealtor(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to get agency of realtor
	agency, err := uh.userService.GetAgencyOfRealtor(ctx, userInfo.ID)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "GET_AGENCY_OF_REALTOR_FAILED", "Failed to get agency of realtor")
		return
	}

	// Convert agency to DTO
	agencyResponse := dto.UserResponse{
		ID:       agency.GetID(),
		FullName: agency.GetFullName(),
		NickName: agency.GetNickName(),
		Email:    agency.GetEmail(),
	}

	// Prepare response
	response := dto.GetAgencyOfRealtorResponse{
		Agency: agencyResponse,
	}

	c.JSON(http.StatusOK, response)
}

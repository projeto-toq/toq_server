package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetRealtorsByAgency(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to get realtors by agency
	realtors, err := uh.userService.GetRealtorsByAgency(ctx, userInfo.ID)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "GET_REALTORS_BY_AGENCY_FAILED", "Failed to get realtors by agency")
		return
	}

	// Convert realtors to DTO
	realtorsResponse := make([]dto.UserResponse, len(realtors))
	for i, realtor := range realtors {
		realtorsResponse[i] = dto.UserResponse{
			ID:          realtor.GetID(),
			FullName:    realtor.GetFullName(),
			NickName:    realtor.GetNickName(),
			NationalID:  realtor.GetNationalID(),
			CreciNumber: realtor.GetCreciNumber(),
			CreciState:  realtor.GetCreciState(),
			Email:       realtor.GetEmail(),
			PhoneNumber: realtor.GetPhoneNumber(),
		}
	}

	// Prepare response
	response := dto.GetRealtorsByAgencyResponse{
		Realtors: realtorsResponse,
	}

	c.JSON(http.StatusOK, response)
}

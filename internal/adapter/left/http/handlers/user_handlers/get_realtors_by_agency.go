package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetRealtorsByAgency(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Get user information from context (set by middleware)
	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok {
		// Se chegar aqui, é erro de pipeline (middleware deveria ter setado)
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	// Call service to get realtors by agency
	realtors, err := uh.userService.GetRealtorsByAgency(ctx, userInfo.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
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

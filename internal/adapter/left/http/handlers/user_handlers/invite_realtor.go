package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) InviteRealtor(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request body
	var request dto.InviteRealtorRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service to invite realtor
	if err := uh.userService.InviteRealtor(ctx, request.PhoneNumber); err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "INVITE_REALTOR_FAILED", "Failed to invite realtor")
		return
	}

	// Prepare response
	response := dto.InviteRealtorResponse{
		Message: "Realtor invitation sent successfully",
	}

	c.JSON(http.StatusOK, response)
}

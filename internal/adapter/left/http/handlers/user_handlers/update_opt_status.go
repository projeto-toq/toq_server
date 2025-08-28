package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) UpdateOptStatus(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request body
	var request dto.UpdateOptStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service to update opt status
	if err := uh.userService.UpdateOptStatus(ctx, request.OptIn); err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "UPDATE_OPT_STATUS_FAILED", "Failed to update opt status")
		return
	}

	// Prepare response
	response := dto.UpdateOptStatusResponse{
		Message: "Opt status updated successfully",
	}

	c.JSON(http.StatusOK, response)
}

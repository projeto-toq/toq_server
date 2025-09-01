package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
)

func (uh *UserHandler) UpdateOptStatus(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request body
	var request dto.UpdateOptStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service to update opt status
	if err := uh.userService.UpdateOptStatus(ctx, request.OptIn); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.UpdateOptStatusResponse{
		Message: "Opt status updated successfully",
	}

	c.JSON(http.StatusOK, response)
}

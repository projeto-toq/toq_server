package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// UpdateOptStatus updates the user's messaging opt-in status
//
//	@Summary      Update opt-in status
//	@Description  Update user's consent to receive notifications (opt-in/opt-out)
//	@Tags         User
//	@Accept       json
//	@Produce      json
//	@Param        request  body      dto.UpdateOptStatusRequest  true  "Opt-in request"
//	@Success      200      {object}  dto.UpdateOptStatusResponse
//	@Failure      400      {object}  dto.ErrorResponse  "Invalid request"
//	@Failure      401      {object}  dto.ErrorResponse  "Unauthorized"
//	@Failure      403      {object}  dto.ErrorResponse  "Forbidden"
//	@Failure      500      {object}  dto.ErrorResponse  "Internal server error"
//	@Router       /user/opt-status [put]
//	@Security     BearerAuth
func (uh *UserHandler) UpdateOptStatus(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

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

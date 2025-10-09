package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// Reference to ensure Swag finds dto.ErrorResponse from annotations
type _ = dto.ErrorResponse

// SignOut handles user sign out
// @Summary		Sign out user
// @Description	Sign out the current user. If refreshToken or deviceToken is provided, revokes a single session/device; otherwise, global signout. When targeting a device without deviceToken, send X-Device-Id.
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			X-Device-Id	header		string	false	"Device ID (UUIDv4) for targeted signout when deviceToken isn't provided"
// @Param			request	body		dto.SignOutRequest	true	"Sign out request"
// @Success		200		{object}	dto.SignOutResponse	"Sign out confirmation message"
// @Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
// @Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
// @Failure		403		{object}	dto.ErrorResponse	"Forbidden"
// @Failure		500		{object}	dto.ErrorResponse	"Internal server error"
// @Router			/user/signout [post]
// @Security		BearerAuth
func (uh *UserHandler) SignOut(c *gin.Context) {
	// Tracing já provido por TelemetryMiddleware.
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Parse request body com DTO consolidado
	var request dto.SignOutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Call service
	if err := uh.userService.SignOut(ctx, request.DeviceToken, request.RefreshToken); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, dto.SignOutResponse{Message: "Successfully signed out"})
}

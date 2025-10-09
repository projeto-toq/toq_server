package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// RequestPhoneChange
//
//	@Summary      Request phone number change
//	@Description  Start a phone change by generating a validation code for the new phone. If a pending change exists (valid or expired), a new code and expiration are generated and persisted, then a notification is sent.
//	@Tags         User
//	@Accept       json
//	@Produce      json
//	@Param        request  body      dto.RequestPhoneChangeRequest  true  "New phone number (E.164)"
//	@Success      200      {object}  dto.RequestPhoneChangeResponse         "Phone change request sent"
//	@Failure      400      {object}  dto.ErrorResponse                      "Invalid request format or phone"
//	@Failure      401      {object}  dto.ErrorResponse                      "Unauthorized"
//	@Failure      409      {object}  dto.ErrorResponse                      "Phone already in use"
//	@Failure      500      {object}  dto.ErrorResponse                      "Internal server error"
//	@Router       /user/phone/request [post]
//	@Security     BearerAuth
func (uh *UserHandler) RequestPhoneChange(c *gin.Context) {
	// Enrich context (request id, user, ip, UA)
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Parse request body
	var request dto.RequestPhoneChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Delegate normalization/validation to the service layer (SSOT reads userID from ctx)
	if err := uh.userService.RequestPhoneChange(ctx, request.NewPhoneNumber); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.RequestPhoneChangeResponse{Message: "Phone change request sent successfully"}
	c.JSON(http.StatusOK, response)
}

package userhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

// RequestEmailChange
//
//	@Summary      Request email change
//	@Description  Start email change by generating a validation code for the new email. If a pending change exists (valid or expired), a new code and expiration are generated and persisted, then a notification is sent.
//	@Tags         User
//	@Accept       json
//	@Produce      json
//	@Param        request  body      dto.RequestEmailChangeRequest  true  "New email"
//	@Success      200      {object}  dto.RequestEmailChangeResponse         "Email change request sent"
//	@Failure      400      {object}  dto.ErrorResponse                      "Invalid request format or email"
//	@Failure      401      {object}  dto.ErrorResponse                      "Unauthorized"
//	@Failure      409      {object}  dto.ErrorResponse                      "Email already in use"
//	@Failure      500      {object}  dto.ErrorResponse                      "Internal server error"
//	@Router       /user/email/request [post]
func (uh *UserHandler) RequestEmailChange(c *gin.Context) {
	// Enriquecer o contexto com informações da requisição e usuário
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Parse request body
	var request dto.RequestEmailChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate and clean email
	email := strings.TrimSpace(request.NewEmail)
	if err := validators.ValidateEmail(email); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_EMAIL", "Invalid email format")
		return
	}

	// Call service to request email change
	if err := uh.userService.RequestEmailChange(ctx, email); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.RequestEmailChangeResponse{
		Message: "Email change request sent successfully",
	}

	c.JSON(http.StatusOK, response)
}

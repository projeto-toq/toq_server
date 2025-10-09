package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
	"github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

// ConfirmEmailChange
//
//	@Summary      Confirm email change
//	@Description  Confirm email change by providing the received validation code
//	@Tags         User
//	@Accept       json
//	@Produce      json
//	@Param        request  body      dto.ConfirmEmailChangeRequest  true  "Confirmation code"
//	@Success      200      {object}  dto.ConfirmEmailChangeResponse         "Confirmation message"
//	@Failure      400      {object}  dto.ErrorResponse                      "Invalid request format or code"
//	@Failure      401      {object}  dto.ErrorResponse                      "Unauthorized"
//	@Failure      409      {object}  dto.ErrorResponse                      "Email change not pending or already in use"
//	@Failure      410      {object}  dto.ErrorResponse                      "Code expired"
//	@Failure      422      {object}  dto.ErrorResponse                      "Invalid code"
//	@Failure      500      {object}  dto.ErrorResponse                      "Internal server error"
//	@Router       /user/email/confirm [post]
func (uh *UserHandler) ConfirmEmailChange(c *gin.Context) {
	// Enriquecer o contexto com informações da requisição e usuário
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Parse request body
	var request dto.ConfirmEmailChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate code
	if err := validators.ValidateCode(request.Code); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_CODE", "Invalid code format")
		return
	}

	// Call service to confirm email change (no tokens returned)
	err := uh.userService.ConfirmEmailChange(ctx, request.Code)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ConfirmEmailChangeResponse{Message: "Email changed successfully"})
}

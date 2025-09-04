package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ResendEmailChangeCode
//
//	@Summary      Resend email change code
//	@Description  Resend the existing (still valid) validation code to the pending new email address. If there is no pending change, returns 409. If the code is expired, returns 410.
//	@Tags         User
//	@Produce      json
//	@Success      200  {object}  dto.ResendEmailChangeCodeResponse  "Confirmation message"
//	@Failure      401  {object}  dto.ErrorResponse                 "Unauthorized"
//	@Failure      409  {object}  dto.ErrorResponse                 "Email change not pending or email already in use"
//	@Failure      410  {object}  dto.ErrorResponse                 "Code expired"
//	@Failure      500  {object}  dto.ErrorResponse                 "Internal server error"
//	@Router       /user/email/resend [post]
func (uh *UserHandler) ResendEmailChangeCode(c *gin.Context) {
	// Enriquecer o contexto com informações da requisição e usuário
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Call service to resend email change code
	if err := uh.userService.ResendEmailChangeCode(ctx); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response (nunca retornar o código no corpo)
	response := dto.ResendEmailChangeCodeResponse{Message: "Code resent to the new email"}
	c.JSON(http.StatusOK, response)
}

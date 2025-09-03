package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ResendPhoneChangeCode
//
//	@Summary      Resend phone change code
//	@Description  Resend a new validation code to the pending new phone number
//	@Tags         User
//	@Produce      json
//	@Success      200  {object}  dto.ResendPhoneChangeCodeResponse  "Confirmation message"
//	@Failure      401  {object}  dto.ErrorResponse                 "Unauthorized"
//	@Failure      409  {object}  dto.ErrorResponse                 "Phone change not pending"
//	@Failure      500  {object}  dto.ErrorResponse                 "Internal server error"
//	@Router       /user/phone/resend [post]
//	@Security     BearerAuth
func (uh *UserHandler) ResendPhoneChangeCode(c *gin.Context) {
	// Enrich context with request info
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Call service to resend phone change code
	if err := uh.userService.ResendPhoneChangeCode(ctx); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response (never return the code in the body)
	response := dto.ResendPhoneChangeCodeResponse{Message: "Code resent to the new phone"}
	c.JSON(http.StatusOK, response)
}

package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteAccount deletes the authenticated user's account
//
//	@Summary      Delete account
//	@Description  Logically delete the current user's account. Revokes all sessions and removes device tokens. User data and role history are preserved for audit purposes and regulatory compliance (LGPD). The account cannot be restored after deletion, but the user may create a new account with the same credentials (email, phone, CPF).
//	@Tags         User
//	@Accept       json
//	@Produce      json
//	@Success      200  {object}  dto.DeleteAccountResponse
//	@Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
//	@Failure      403  {object}  dto.ErrorResponse  "Forbidden"
//	@Failure      500  {object}  dto.ErrorResponse  "Internal server error"
//	@Router       /user/account [delete]
//	@Security     BearerAuth
func (uh *UserHandler) DeleteAccount(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	// Call service to delete account
	tokens, err := uh.userService.DeleteAccount(ctx)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.DeleteAccountResponse{
		Tokens: dto.TokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
		Message: "Account successfully deleted",
	}

	c.JSON(http.StatusOK, response)
}

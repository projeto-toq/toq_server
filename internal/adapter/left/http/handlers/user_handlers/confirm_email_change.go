package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

// ConfirmEmailChange
//
//	@Summary      Confirm email change
//	@Description  Confirm email change by providing the received validation code
//	@Tags         User
//	@Accept       json
//	@Produce      json
//	@Param        request  body      dto.ConfirmEmailChangeRequest  true  "Confirmation code"
//	@Success      200      {object}  dto.ConfirmEmailChangeResponse         "Tokens returned if applicable"
//	@Failure      400      {object}  dto.ErrorResponse                      "Invalid request format or code"
//	@Failure      401      {object}  dto.ErrorResponse                      "Unauthorized"
//	@Failure      409      {object}  dto.ErrorResponse                      "Email change not pending or already in use"
//	@Failure      410      {object}  dto.ErrorResponse                      "Code expired"
//	@Failure      422      {object}  dto.ErrorResponse                      "Invalid code"
//	@Failure      500      {object}  dto.ErrorResponse                      "Internal server error"
//	@Router       /user/email/change/confirm [post]
func (uh *UserHandler) ConfirmEmailChange(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

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

	// Call service to confirm email change
	tokens, err := uh.userService.ConfirmEmailChange(ctx, userInfo.ID, request.Code)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.ConfirmEmailChangeResponse{
		Tokens: dto.TokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}

	c.JSON(http.StatusOK, response)
}

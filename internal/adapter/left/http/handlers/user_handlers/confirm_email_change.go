package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

func (uh *UserHandler) ConfirmEmailChange(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Parse request body
	var request dto.ConfirmEmailChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate code
	if err := validators.ValidateCode(request.Code); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_CODE", "Invalid code format")
		return
	}

	// Call service to confirm email change
	tokens, err := uh.userService.ConfirmEmailChange(ctx, userInfo.ID, request.Code)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "CONFIRM_EMAIL_CHANGE_FAILED", "Failed to confirm email change")
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

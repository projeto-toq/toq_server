package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) DeleteAccount(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to delete account
	tokens, err := uh.userService.DeleteAccount(ctx, userInfo.ID)
	if err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "DELETE_ACCOUNT_FAILED", "Failed to delete account")
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

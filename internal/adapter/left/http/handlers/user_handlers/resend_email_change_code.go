package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func (uh *UserHandler) ResendEmailChangeCode(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to resend email change code
	code, err := uh.userService.ResendEmailChangeCode(ctx, userInfo.ID)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "RESEND_EMAIL_CHANGE_CODE_FAILED", "Failed to resend email change code")
		return
	}

	// Prepare response
	response := dto.ResendEmailChangeCodeResponse{
		Code: code,
	}

	c.JSON(http.StatusOK, response)
}

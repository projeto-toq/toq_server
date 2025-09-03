package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// ResendEmailChangeCode
//
//	@Summary      Resend email change code
//	@Description  Resend a new validation code to the pending new email address
//	@Tags         User
//	@Produce      json
//	@Success      200  {object}  dto.ResendEmailChangeCodeResponse  "Confirmation message"
//	@Failure      401  {object}  dto.ErrorResponse                 "Unauthorized"
//	@Failure      409  {object}  dto.ErrorResponse                 "Email change not pending"
//	@Failure      500  {object}  dto.ErrorResponse                 "Internal server error"
//	@Router       /user/email/change/resend [post]
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
	if err := uh.userService.ResendEmailChangeCode(ctx, userInfo.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response (nunca retornar o código no corpo)
	response := dto.ResendEmailChangeCodeResponse{Message: "Code resent to the new email"}
	c.JSON(http.StatusOK, response)
}

package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
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
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Call service to resend phone change code
	if err := uh.userService.ResendPhoneChangeCode(ctx, userInfo.ID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response (never return the code in the body)
	response := dto.ResendPhoneChangeCodeResponse{Message: "Code resent to the new phone"}
	c.JSON(http.StatusOK, response)
}

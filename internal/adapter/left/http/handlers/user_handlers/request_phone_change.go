package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// RequestPhoneChange
//
//	@Summary      Request phone number change
//	@Description  Start a phone change by sending a code to the new phone number
//	@Tags         User
//	@Accept       json
//	@Produce      json
//	@Param        request  body      dto.RequestPhoneChangeRequest  true  "New phone number (E.164)"
//	@Success      200      {object}  dto.RequestPhoneChangeResponse         "Phone change request sent"
//	@Failure      400      {object}  dto.ErrorResponse                      "Invalid request format or phone"
//	@Failure      401      {object}  dto.ErrorResponse                      "Unauthorized"
//	@Failure      409      {object}  dto.ErrorResponse                      "Phone already in use or same as current"
//	@Failure      500      {object}  dto.ErrorResponse                      "Internal server error"
//	@Router       /user/phone/request [post]
//	@Security     BearerAuth
func (uh *UserHandler) RequestPhoneChange(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Parse request body
	var request dto.RequestPhoneChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Delegate normalization/validation to the service layer
	if err := uh.userService.RequestPhoneChange(ctx, userInfo.ID, request.NewPhoneNumber); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.RequestPhoneChangeResponse{Message: "Phone change request sent successfully"}
	c.JSON(http.StatusOK, response)
}

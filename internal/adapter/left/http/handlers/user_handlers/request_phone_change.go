package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/converters"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

func (uh *UserHandler) RequestPhoneChange(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Parse request body
	var request dto.RequestPhoneChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate and clean phone number
	newPhone := converters.RemoveAllButDigitsAndPlusSign(request.NewPhoneNumber)
	if err := validators.ValidateE164(newPhone); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_PHONE", "Invalid phone number format")
		return
	}

	// Call service to request phone change
	if err := uh.userService.RequestPhoneChange(ctx, userInfo.ID, newPhone); err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "REQUEST_PHONE_CHANGE_FAILED", "Failed to request phone change")
		return
	}

	// Prepare response
	response := dto.RequestPhoneChangeResponse{
		Message: "Phone change request sent successfully",
	}

	c.JSON(http.StatusOK, response)
}
